package cron

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
)

func validateName(name string) error {
	// if name == "" {
	// 	return status.Error(codes.InvalidArgument, "cron job name cannot be empty")
	// }
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
	data, err := readCrontab(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "read crontab: %s", err.Error())
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20) // default 64 KB line limit; set it to 1MB
	for sc.Scan() {
		line := sc.Text()
		schedule, command, err := ParseLine(line)
		if err != nil {
			log.Printf("failed to parse cron line: %s", err)
			continue
		}
		hash := hashCronJob(schedule, command)
		job, err := db.FindCronJobByHash(database, hash)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newCronJob := &db.CronJob{
					Hash:    hash,
					Name:    "",
					Status:  db.CronJobStatusScheduled,
					Managed: strings.HasSuffix(command, fmt.Sprintf("%s run-job", common.GetAgentPath())),
				}
				if err := database.Create(newCronJob).Error; err != nil {
					log.Printf("failed to create cron job: %s", err)
					return status.Errorf(codes.Internal, "create cron job: %s", err.Error())
				}
				job = newCronJob

			} else {
				log.Printf("failed to find cron job by hash: %s", err)
				return status.Errorf(codes.Internal, "find cron job: %s", err.Error())
			}
		}
		// update managed status
		job.Managed = strings.HasSuffix(command, fmt.Sprintf("%s run-job", common.GetAgentPath()))
		if err := database.Save(job).Error; err != nil {
			log.Printf("failed to update cron job: %s", err)
			return status.Errorf(codes.Internal, "update cron job: %s", err.Error())
		}
		if err := onLine(job, schedule, command); err != nil {
			return err
		}
	}

	return nil
}

func readCrontab(ctx context.Context) ([]byte, error) {
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
	return out.Bytes(), nil
}

func hashCronJob(schedule, command string) string {
	h := sha256.New()
	for _, s := range []string{
		strings.Join(strings.Fields(schedule), " "), // normalize whitespace
		strings.TrimSpace(command),
	} {
		fmt.Fprintf(h, "%d:%s", len(s), s) // length prefix, no ambiguity
	}
	return hex.EncodeToString(h.Sum(nil))
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
