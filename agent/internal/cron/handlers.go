package cron

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListCronJobs(database *gorm.DB, req *proto.ListCronJobsRequest, stream proto.CronService_ListCronJobsServer) error {
	hashes := make([]string, 0)

	batch := make([]*proto.Cron, 0, 500)
	err := readCronJobEach(stream.Context(), database, func(job *db.CronJob, schedule, command string) error {
		if len(batch) == 500 {
			stream.Send(&proto.ListCronJobsResponse{Jobs: batch})

			batch = batch[:0] // reset batch
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
		log.Error(err)
		return status.Errorf(codes.Internal, "failed to list cron jobs: %v", err)
	}
	go func() {
		if len(hashes) > 0 {
			// delete orphaned jobs
			err := database.Where("hash NOT IN ?", hashes).Delete(&db.CronJob{}).Error
			if err != nil {
				log.Printf("failed to delete orphaned cron jobs: %s", err)
			}
		}
	}()
	stream.Send(&proto.ListCronJobsResponse{Jobs: batch, Done: true})

	return nil
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
			command = fmt.Sprintf("%s run %d", common.GetAgentPath(), newCronJob.ID)
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
