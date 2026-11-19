package cron

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
)

func GetCronJobCommandByHash(ctx context.Context, hash string) (string, error) {
	data, err := readCrontab(ctx)
	if err != nil {
		return "", err
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20) // default 64 KB line limit; set it to 1MB
	for sc.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		line := sc.Text()
		schedule, command, err := ParseLine(line)
		if err != nil {
			slog.Error("failed to parse cron line", "error", err)
			continue
		}
		h := hashCronJob(schedule, command)
		if h == hash {
			return command, nil
		}
	}

	return "", fmt.Errorf("Could not find command for cron job with hash %q", hash)
}
func validateName(name string) error {

	if name == "." || name == ".." {
		return status.Errorf(codes.InvalidArgument, "invalid cron job name %q", name)
	}
	if len(name) > 255 {
		return status.Errorf(codes.InvalidArgument, "cron job name %q is too long", name)
	}
	if strings.ContainsAny(name, `/\`) || strings.ContainsRune(name, '\x00') || strings.ContainsRune(name, ':') || strings.ContainsRune(name, '*') || strings.ContainsRune(name, '?') || strings.ContainsRune(name, '"') || strings.ContainsRune(name, '<') || strings.ContainsRune(name, '>') || strings.ContainsRune(name, '|') || strings.Contains(name, "..") {
		return status.Errorf(codes.InvalidArgument, "invalid cron job name %q", name)
	}

	return nil
}

// Read and process each cron job from the crontab.
func readCronJobEach(ctx context.Context, database *gorm.DB, onLine func(job *db.CronJob, schedule, command string) error) error {
	slog.Debug("reading crontab for processing")
	data, err := readCrontab(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "read crontab: %s", err.Error())
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20) // default 64 KB line limit; set it to 1MB
	slog.Debug("Now scanning crontab lines", "data_size", len(data))
	for sc.Scan() {
		line := sc.Text()
		schedule, command, err := ParseLine(line)
		if err != nil {
			slog.Error("failed to parse cron line", "error", err)
			continue
		}
		hash := hashCronJob(schedule, command)
		job, err := db.FindCronJobByHash(database, hash)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newCronJob := &db.CronJob{
					Hash:    hash,
					Name:    "",
					Status:  db.CronJobStatusActive,
					Managed: false,
				}
				if err := database.Create(newCronJob).Error; err != nil {
					slog.Error("failed to create cron job", "error", err)
					return status.Errorf(codes.Internal, "create cron job: %s", err.Error())
				}
				job = newCronJob

			} else {
				slog.Error("failed to find cron job by hash", "error", err)
				return status.Errorf(codes.Internal, "find cron job: %s", err.Error())
			}
		}
		// update managed status
		job.Managed = strings.HasPrefix(command, fmt.Sprintf("%s run-job", common.GetAgentPath()))
		command, id := splitID(trimAgentPrefix(command))
		if job.Managed {
			if id != job.ID {
				slog.Info("cron job ID mismatch, updating crontab", "job_id", job.ID, "line_id", id)
				// update crontab with new ID
				if err := updateJob(ctx, job.Hash, schedule, addIDAndAgentPrefix(command, job.ID)); err != nil {
					slog.Error("failed to update crontab", "error", err)
					return status.Errorf(codes.Internal, "update crontab: %s", err.Error())
				}
			}
		}
		nextRunAt, err := NextRun(schedule, time.Now())
		if err != nil {
			slog.Error("failed to calculate next run time", "error", err)
			continue
		}
		job.NextRunAt = nextRunAt

		if err := database.Save(job).Error; err != nil {
			slog.Error("failed to update cron job", "error", err)
			return status.Errorf(codes.Internal, "update cron job: %s", err.Error())
		}
		if err := onLine(job, schedule, command); err != nil {
			return err
		}
	}

	return nil
}

func updateJob(ctx context.Context, hash, newSchedule, newCommand string) error {
	slog.Debug("updating cron job", "hash", hash)
	data, err := readCrontab(ctx)
	if err != nil {
		return fmt.Errorf("read crontab: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	removed := false
	newLine, err := BuildLine(newSchedule, newCommand)
	if err != nil {
		return fmt.Errorf("build line: %w", err)
	}
	slog.Debug("Now looping cron jobs in update, ", "lines", len(lines))
	for _, line := range lines {
		schedule, command, err := ParseLine(line)
		if err != nil {
			slog.Error("failed to parse cron line", "error", err)
			// error, skip parsing but preserve line
			out = append(out, line)
			continue
		}

		if hashCronJob(schedule, command) == hash {
			removed = true
			out = append(out, newLine)
			continue
		}
		// add lines
		out = append(out, line)
	}

	if !removed {
		slog.Info("cron job not found, adding new.", "hash", hash)
		out = append(out, newLine)
	}

	content := []byte(strings.Join(out, "\n"))
	if err := installCrontab(ctx, content); err != nil {
		return fmt.Errorf("install crontab: %w", err)
	}
	return nil
}

func readCrontab(ctx context.Context) ([]byte, error) {
	slog.Debug("reading crontab")
	var out, errb bytes.Buffer
	cmd := exec.CommandContext(ctx, "crontab", "-l")
	cmd.Stdout = &out
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ExitCode() == 1 &&
			strings.Contains(errb.String(), "no crontab for") {
			return []byte{}, nil // empty crontab is not an error
		}
		return nil, fmt.Errorf("crontab -l: %w: %s", err, strings.TrimSpace(errb.String()))
	}
	slog.Debug("read crontab successfully")
	return out.Bytes(), nil
}
func addAgentPrefix(command string) string {
	prefix := getAgentPrefix()
	return fmt.Sprintf("%s %s", prefix, command)
}
func addIDAndAgentPrefix(command string, id uint64) string {
	return addAgentPrefix(addID(command, id))
}
func addID(command string, id uint64) string {
	return fmt.Sprintf("%s id=%d", command, id)
}
func getAgentPrefix() string {
	return fmt.Sprintf("%s run-job", common.GetAgentPath())
}
func trimAgentPrefix(command string) string {
	prefix := getAgentPrefix()
	return strings.TrimPrefix(command, prefix)
}
func hashCronJob(schedule, command string) string {
	slog.Debug("hashing cron job", "schedule", schedule)
	// remove agent prefix and trailing ID
	command = trimAgentPrefix(command)
	command, _ = splitID(command)
	h := sha256.New()
	for _, s := range []string{
		strings.Join(strings.Fields(schedule), " "), // normalize whitespace
		strings.TrimSpace(command),
	} {
		fmt.Fprintf(h, "%d:%s", len(s), s) // length prefix, no ambiguity
	}

	return hex.EncodeToString(h.Sum(nil))
}

// dedupeCronLines removes lines that describe the same job (same hash),
// keeping the first occurrence. Non-job lines (comments, blanks, env
// assignments) are preserved.
func dedupeCronLines(data []byte) ([]byte, uint32) {
	lines := strings.Split(string(data), "\n")
	seen := make(map[string]struct{}, len(lines))
	out := make([]string, 0, len(lines))
	var removed uint32

	for _, line := range lines {
		schedule, command, err := ParseLine(line)
		if err != nil {
			out = append(out, line)
			continue
		}

		hash := hashCronJob(schedule, command)
		if _, exists := seen[hash]; exists {
			removed++
			continue
		}

		seen[hash] = struct{}{}
		out = append(out, line)
	}

	return []byte(strings.Join(out, "\n")), removed
}

// installCrontab replaces the user's crontab with content, via stdin.
func installCrontab(ctx context.Context, content []byte) error {
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "crontab", "-")
	cmd.Stdin = bytes.NewReader(content)
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("crontab install: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func appendLine(cur []byte, line string) []byte {
	out := append([]byte(nil), cur...)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	return append(out, line+"\n"...)
}

func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func splitID(s string) (string, uint64) {
	before, after, ok := strings.CutLast(s, " id=")
	if !ok {
		return strings.TrimSpace(s), 0
	}
	id, err := strconv.ParseUint(strings.TrimSpace(after), 10, 64)
	if err != nil {
		return strings.TrimSpace(s), 0 // " id=" wasn't a real trailing ID, leave the text alone
	}
	return strings.TrimSpace(before), id
}
