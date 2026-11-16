package models

import "hazyio.com/servergram/shared/gen/proto"

type CreateCronJobRequest struct {
	ID   uint `path:"id" minimum:"1" doc:"The ID of the current resource"`
	Body struct {
		Schedule          string `json:"schedule" validate:"required"`
		Command           string `json:"command" validate:"required"`
		Name              string `json:"name"`
		MetaKeepCount     int32  `json:"meta_keep_count" validate:"gte=0,lte=960"`
		RunLogsKeepCount  int32  `json:"run_logs_keep_count" validate:"gte=0,lte=960"`
		ErrorLogsKeepDays int32  `json:"error_logs_keep_days" validate:"gte=0,lte=365"`
		Wrap              bool   `json:"wrap"`
	}
}

type CreateCronJobResponse struct {
	Body *proto.Cron
}

type RemoveDuplicateCronJobsResponse struct {
	Body struct {
		Removed uint32 `json:"removed" doc:"Number of duplicate cron job lines removed"`
	}
}
