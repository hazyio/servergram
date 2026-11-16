package cron

import (
	"time"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/db"
	"hazyio.com/servergram/shared/gen/proto"
)

func handleListCronJobs(database *gorm.DB, req *proto.ListCronJobsRequest, stream proto.CronService_ListCronJobsServer) error {

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
		return nil
	})
	if err != nil {
		log.Error(err)
		return status.Errorf(codes.Internal, "failed to list cron jobs: %v", err)
	}

	stream.Send(&proto.ListCronJobsResponse{Jobs: batch, Done: true})

	return nil
}

func handleCreateCronJob(database *gorm.DB, req *proto.CreateCronJobRequest) (*proto.CreateCronJobResponse, error) {
	if err := validateName(req.Name); err != nil {
		return nil, err
	}
	if err := ValidateCommand(req.Command); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid command: %v", err)
	}
	
	if err := ValidateSchedule(req.Schedule); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid schedule: %v", err)
	}
	
}
