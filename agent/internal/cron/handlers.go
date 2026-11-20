package cron

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
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
	cronJobHash := hashCronJob(req.Schedule, req.Name)
	var newCronJobLine string
	newCronJob := db.CronJob{
		Name:             req.GetName(),
		Hash:             cronJobHash,
		Managed:          req.GetWrap(),
		Status:           db.CronJobStatusActive,
		MetaKeepCount:    req.GetMetaKeepCount(),
		ErrorLogKeepDays: req.GetErrorLogsKeepDays(),
		RunLogKeepCount:  req.GetRunLogsKeepCount(),
	}
	bound1 := []int{5, 10, 20, 40, 80, 160, 240, 480, 960}
	errorDaysBounds := []int{1, 7, 14, 30, 60, 90, 180, 365}
	if req.Wrap {
		if !slices.Contains(bound1, int(req.GetMetaKeepCount())) {
			return nil, status.Errorf(codes.InvalidArgument, "meta_keep_count must be one of 5, 10, 20, 40, 80, 160, 240, 480, or 960")
		}
		if !slices.Contains(bound1, int(req.GetRunLogsKeepCount())) {
			return nil, status.Errorf(codes.InvalidArgument, "run_logs_keep_count must be one of 5, 10, 20, 40, 80, 160, 240, 480, or 960")
		}
		if !slices.Contains(errorDaysBounds, int(req.GetErrorLogsKeepDays())) {
			return nil, status.Errorf(codes.InvalidArgument, "error_logs_keep_days must be one of 1, 7, 14, 30, 60, 90, 180, or 365")
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
	cmd := exec.Command(common.GetAgentPath(), "run-job", jobCommand, fmt.Sprintf("id=%d", job.ID))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // new session, no controlling terminal

	if err := cmd.Start(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start cron job: %v", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release cron job process: %v", err)
	}
	return &proto.RunCronJobResponse{}, nil
}

func handleStopCronJob(ctx context.Context, database *gorm.DB, req *proto.StopCronJobRequest) (*proto.StopCronJobResponse, error) {
	job, err := db.FindCronJobByID(database, uint64(req.Id))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found: %v", err)
	}
	if !job.Managed {
		return nil, status.Errorf(codes.FailedPrecondition, "cron job is not managed by agent")
	}
	job.Status = db.CronJobStatusInactive
	if err := database.Save(job).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update cron job status: %v", err)
	}
	return &proto.StopCronJobResponse{}, nil
}
func handleInfo(ctx context.Context, database *gorm.DB, req *proto.CronInfoRequest) (*proto.CronInfoResponse, error) {
	var job *proto.CronInfoResponse
	var FoundErr = errors.New("not found")

	err := readCronJobEach(ctx, database, func(j *db.CronJob, schedule, command string) error {
		if j.ID != req.Id {
			return nil
		}
		job = &proto.CronInfoResponse{
			Name:              j.Name,
			Managed:           j.Managed,
			Id:                j.ID,
			Schedule:          schedule,
			Command:           command,
			MetaKeepCount:     j.MetaKeepCount,
			RunLogsKeepCount:  j.RunLogKeepCount,
			ErrorLogsKeepDays: j.ErrorLogKeepDays,
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
		job.MetaKeepCount = req.MetaKeepCount
		job.RunLogKeepCount = req.RunLogsKeepCount
		job.ErrorLogKeepDays = req.ErrorLogsKeepDays
		if err := tx.Save(job).Error; err != nil {
			return err
		}
		if job.Managed {
			command = addIDAndAgentPrefix(command, job.ID)
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
	err = database.Transaction(func(tx *gorm.DB) error {
		var file []byte
		found := false

		readErr := readCronJobEach(ctx, tx, func(curJob *db.CronJob, schedule, command string) error {
			if curJob.Hash == job.Hash {
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
			return readErr
		}
		if !found {
			return fmt.Errorf("cron job not found in crontab")
		}
		if err := database.Delete(job).Error; err != nil {
			return err
		}
		if err := installCrontab(ctx, file); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error("failed to delete cron job", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete cron job: %v", err)
	}

	return &proto.DeleteCronJobResponse{}, nil
}

func handleDuplicateCronJob(ctx context.Context, database *gorm.DB, req *proto.DuplicateCronJobRequest) (*proto.DuplicateCronJobResponse, error) {
	job, err := db.FindCronJobByID(database, uint64(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}
	var file []byte
	found := false
	duplicated := false
	err = readCronJobEach(ctx, database, func(curJob *db.CronJob, schedule, command string) error {
		if curJob.Hash == job.Hash {
			// found target job,
			found = true
		}
		if curJob.Managed {
			command = addIDAndAgentPrefix(command, curJob.ID)
		}
		line, err := BuildLine(schedule, command)
		if err != nil {
			return err
		}
		file = appendLine(file, line)
		if found && !duplicated {
			duplicated = true
			file = appendLine(file, line)

		}

		return nil
	})
	if err != nil {
		slog.Error("failed to read cron jobs", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to duplicate cron job: %v", err)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "cron job not found")
	}

	if err := installCrontab(ctx, file); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to install crontab: %v", err)
	}
	return &proto.DuplicateCronJobResponse{}, nil
}
