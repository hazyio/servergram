package cron

import (
	"time"

	"github.com/labstack/gommon/log"
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
		return err
	}

	stream.Send(&proto.ListCronJobsResponse{Jobs: batch, Done: true})

	return nil
}
