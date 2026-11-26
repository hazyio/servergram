package cron

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListCronJobs(database *gorm.DB, stream proto.CronService_ListCronJobsServer) error {
	slog.Debug("listing cron jobs")
	hashes := make([]string, 0)
	ids := make([]uint64, 0)
	batch := make([]*proto.Cron, 0, 50)
	err := readCronJobEach(stream.Context(), database, func(job *db.CronJob, schedule, command string) error {
		if len(batch) == 50 {
			slog.Debug("sending batch of 50 cron jobs")
			if err := common.LogError(stream.Send(&proto.ListCronJobsResponse{Jobs: batch})); err != nil {
				return err
			}
			batch = make([]*proto.Cron, 0, 50) // new slice, never reuse after Send
		}

		batch = append(batch, &proto.Cron{
			Hash:      job.Hash,
			Name:      job.Name,
			Managed:   job.Managed,
			Status:    string(job.Status),
			Id:        job.ID,
			LastRunAt: job.LastRunAt.Format(time.RFC3339),
			NextRunAt: job.NextRunAt.Format(time.RFC3339),
			Schedule:  schedule,
			Command:   command,
		})
		hashes = append(hashes, job.Hash)
		ids = append(ids, job.ID)

		return nil
	})
	if err != nil {
		slog.Error("failed to read cron jobs", "error", err)
		return status.Errorf(codes.Internal, "failed to list cron jobs: %v", err)
	}
	go func() {
		if len(hashes) > 0 {
			// delete orphaned jobs
			err := database.Where("hash NOT IN ?", hashes).Delete(&db.CronJob{}).Error
			if err != nil {
				slog.Error("failed to delete orphaned cron jobs", "error", err)
			}

		}
		if len(ids) > 0 {
			// remove orphaned logs
			pruneOrpahanedCronRunsLogs(ids)
		}
	}()
	slog.Debug("sending final batch of cron jobs", "jobs", len(batch))
	if err := common.LogError(stream.Send(&proto.ListCronJobsResponse{Jobs: batch, Done: true})); err != nil {
		return err
	}

	return nil
}

func handleRemoveDuplicateCronJobs(ctx context.Context) (*proto.RemoveDuplicateCronJobsResponse, error) {
	data, err := readCrontab(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "read crontab: %v", err)
	}

	cleaned, removed := dedupeCronLines(data)
	if removed > 0 {
		if err := installCrontab(ctx, cleaned); err != nil {
			return nil, status.Errorf(codes.Internal, "install crontab: %v", err)
		}
	}

	return &proto.RemoveDuplicateCronJobsResponse{Removed: removed}, nil
}

func handleCreateCronJob(ctx context.Context, database *gorm.DB, req *proto.CreateCronJobRequest) (*proto.CreateCronJobResponse, error) {
	if err := validateName(req.Name); err != nil {
		return nil, err
	}
	if err := ValidateCommand(req.Command); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid command: %v", err)
	}

	if err := ValidateSchedule(req.Schedule); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid schedule: %v", err)
	}
	cronJobHash := ""
	if req.Wrap {
		cronJobHash = hashCronJob(req.Schedule, addIDAndAgentPrefix(shQuote(req.Command), 0))

	} else {
		cronJobHash = hashCronJob(req.Schedule, req.Command)
	}
	var newCronJobLine string
	newCronJob := db.CronJob{
		Name:            req.Name,
		Hash:            cronJobHash,
		Managed:         req.Wrap,
		Status:          db.CronJobStatusActive,
		RunLogKeepCount: req.RunLogsKeepCount,
	}
	if req.Wrap {

		if !slices.Contains([]int{5, 10, 20, 40, 80, 160, 240, 480, 960}, int(newCronJob.RunLogKeepCount)) {
			return nil, status.Errorf(codes.InvalidArgument, "run_logs_keep_count must be one of 5, 10, 20, 40, 80, 160, 240, 480, or 960")
		}

	}
	var ErrDuplicateJob = errors.New("a cron job with this schedule and command already exists")

	txErr := database.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&newCronJob).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDuplicateJob
			}
			return err
		}
		var command string
		if req.Wrap {
			command = addIDAndAgentPrefix(shQuote(req.Command), newCronJob.ID)
		} else {
			command = fmt.Sprintf("%s", req.Command)
		}
		newCronJobLine, err = BuildLine(req.Schedule, command)
		if err != nil {
			return err
		}
		cur, err := readCrontab(ctx)
		if err != nil {
			return err
		}
		return installCrontab(ctx, appendLine(cur, newCronJobLine))

	})
	if txErr != nil {
		return nil, status.Errorf(codes.Internal, "failed to create new cron job: %v", txErr)
	}
	return &proto.CreateCronJobResponse{}, nil

}

func handleFixUnmanagedCronJobs(ctx context.Context, database *gorm.DB) (*proto.FixUnmanagedCronJobsResponse, error) {
	count := uint32(0)
	file := make([]byte, 0)
	err := readCronJobEach(ctx, database, func(job *db.CronJob, schedule, command string) error {
		if !job.Managed {
			count++
		}
		line, err := BuildLine(schedule, addIDAndAgentPrefix(command, job.ID))
		if err != nil {
			return err
		}
		file = appendLine(file, line)

		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cron jobs: %v", err)
	}
	if count > 0 {
		if err := installCrontab(ctx, file); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to install crontab: %v", err)
		}
	}
	return &proto.FixUnmanagedCronJobsResponse{Fixed: count}, nil
}

func handleRunCronJob(ctx context.Context, database *gorm.DB, req *proto.RunCronJobRequest) (*proto.RunCronJobResponse, error) {
	job, err := db.FindCronJobByID(database, uint64(req.Id))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found: %v", err)
	}

	jobCommand, err := GetCronJobCommandByHash(ctx, job.Hash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job missing, in database but not in crontab: %v", err)
	}
	slog.Debug("running cron job", "job_id", job.ID, "command", jobCommand)
	parts := strings.Fields(jobCommand) // handles repeated spaces, unlike Split
	if len(parts) == 0 {
		return nil, status.Errorf(codes.Internal, "empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // new session, no controlling terminal

	if err := cmd.Start(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start cron job: %v", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release cron job process: %v", err)
	}
	return &proto.RunCronJobResponse{}, nil
}

func handleStopCronJobRun(ctx context.Context, database *gorm.DB, req *proto.StopCronJobRunRequest) (*proto.StopCronJobRunResponse, error) {
	job, err := db.FindCronJobRunByID(database, uint64(req.Id))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job run not found: %v", err)
	}
	if job.Status != db.RunStatusRunning {
		return nil, status.Errorf(codes.InvalidArgument, "cannot stop cron job run with status: %s", job.Status)
	}

	job.Status = db.RunStatusCancelled
	if err := database.Save(job).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update cron job status: %v", err)
	}
	return &proto.StopCronJobRunResponse{}, nil
}
func handleInfo(ctx context.Context, database *gorm.DB, req *proto.CronInfoRequest) (*proto.CronInfoResponse, error) {
	var job *proto.CronInfoResponse
	var FoundErr = errors.New("not found")

	err := readCronJobEach(ctx, database, func(j *db.CronJob, schedule, command string) error {
		if j.ID != req.Id {
			return nil
		}
		job = &proto.CronInfoResponse{
			Status:           string(j.Status),
			Name:             j.Name,
			Managed:          j.Managed,
			Id:               j.ID,
			Schedule:         schedule,
			Command:          command,
			RunLogsKeepCount: j.RunLogKeepCount,
			LastRunAt:        j.LastRunAt.Format(time.RFC3339),
			NextRunAt:        j.NextRunAt.Format(time.RFC3339),
		}
		return FoundErr
	})

	if err != nil && !errors.Is(err, FoundErr) {
		slog.Error("failed to calculate next run time", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to calculate next run time")
	}
	if job == nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}

	return job, nil
}
func handleUpdateCronJob(ctx context.Context, database *gorm.DB, req *proto.UpdateCronJobRequest) (*proto.UpdateCronJobResponse, error) {
	job, err := db.FindCronJobByID(database, uint64(req.Id))

	if err != nil {
		slog.Error("failed to find cron job", "error", err)
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}
	schedule, command, err := ParseLine(fmt.Sprintf("%s %s", req.Schedule, req.Command))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid schedule or command")
	}
	if err := validateName(req.Name); err != nil {
		return nil, err
	}
	err = database.Transaction(func(tx *gorm.DB) error {
		job.Name = req.Name
		job.RunLogKeepCount = req.RunLogsKeepCount
		if req.Wrap {
			command = addIDAndAgentPrefix(command, job.ID)
		}
		hashCommand := hashCronJob(schedule, command)
		job.Hash = hashCommand
		job.Managed = req.Wrap
		if err := tx.Save(job).Error; err != nil {
			return err
		}

		if err := updateJob(ctx, job.Hash, schedule, command); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		slog.Error("failed to update cron job", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update cron job")
	}

	return &proto.UpdateCronJobResponse{}, nil
}

func handleDeleteCronJob(ctx context.Context, database *gorm.DB, req *proto.DeleteCronJobRequest) (*proto.DeleteCronJobResponse, error) {
	job, err := db.FindCronJobByID(database, uint64(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}
	var file []byte
	found := false
	readErr := readCronJobEach(ctx, database, func(curJob *db.CronJob, schedule, command string) error {
		if curJob.Hash == job.Hash && !found {
			// found target job, skip it
			found = true
			return nil
		}
		if curJob.Managed {
			command = addIDAndAgentPrefix(command, curJob.ID)
		}
		line, err := BuildLine(schedule, command)
		if err != nil {
			return err
		}
		file = appendLine(file, line)

		return nil
	})
	if readErr != nil {
		return nil, status.Errorf(codes.Internal, "An error occurred reading cron job: %v", readErr)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "cron job not found in crontab")
	}

	if err := installCrontab(ctx, file); err != nil {
		slog.Error("An error occurred deleting cron job", "error", err)

		return nil, status.Errorf(codes.Internal, "An error occurred deleting cron job")
	}
	if err := database.Delete(job).Error; err != nil {
		slog.Error("An error occured deleting cron job", "error", err)
		return nil, status.Errorf(codes.Internal, "An error occurred deleting cron job")
	}
	return &proto.DeleteCronJobResponse{}, nil
}

func handleListCronJobRuns(database *gorm.DB, req *proto.ListCronJobRunsRequest) (*proto.ListCronJobRunsResponse, error) {
	_, err := db.FindCronJobByID(database, uint64(req.CronJobId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}

	runs := make([]*db.CronRun, 0, req.Limit)
	err = database.Where("cron_job_id = ?", uint64(req.CronJobId)).Order("id desc").Offset(int(req.Offset)).Limit(int(req.Limit)).Find(&runs).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Internal, "failed to get run meta list: %v", err)
		}
		return &proto.ListCronJobRunsResponse{Runs: make([]*proto.CronRunMeta, 0), More: false}, nil
	}
	thereIsMore := false
	if len(runs) > 0 {
		lastId := runs[len(runs)-1].ID
		var one *proto.CronRunMeta
		if err := database.Where("cron_job_id = ? AND id > ?", uint64(req.CronJobId), lastId).First(&one).Error; err == nil {
			thereIsMore = true
		}
	}
	return &proto.ListCronJobRunsResponse{Runs: common.Map(runs, func(run *db.CronRun) *proto.CronRunMeta {
		return &proto.CronRunMeta{
			Id:         uint32(run.ID),
			Status:     string(run.Status),
			StartedAt:  run.StartedAt.Format(time.RFC3339),
			FinishedAt: run.FinishedAt.Format(time.RFC3339),
			ExitCode:   int32(run.ExitCode),
			Trigger:    string(run.Trigger),
		}
	}), More: thereIsMore}, nil
}
