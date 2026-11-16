package cron

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListCronJobs(database *gorm.DB, stream proto.CronService_ListCronJobsServer) error {
	hashes := make([]string, 0)

	batch := make([]*proto.Cron, 0, 50)
	err := readCronJobEach(stream.Context(), database, func(job *db.CronJob, schedule, command string) error {
		if len(batch) == 50 {
			slog.Debug("sending batch of cron jobs")
			if err := common.LogError(stream.Send(&proto.ListCronJobsResponse{Jobs: batch})); err != nil {
				return err
			}
			batch = make([]*proto.Cron, 0, 50) // new slice, never reuse after Send
		}
		nextRunAt, err := NextRun(schedule, time.Now())
		if err != nil {
			return err
		}
		batch = append(batch, &proto.Cron{
			Hash:      job.Hash,
			Name:      job.Name,
			Managed:   job.Managed,
			Status:    string(job.Status),
			Id:        job.ID,
			LastRunAt: job.LastRunAt.Format(time.RFC3339),
			NextRunAt: nextRunAt.Format(time.RFC3339),
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
		Status:           db.CronJobStatusScheduled,
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
