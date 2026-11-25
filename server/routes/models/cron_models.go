package models

import "hazyio.com/servergram/shared/gen/proto"

type CreateCronJobRequest struct {
	ID   uint `path:"id" minimum:"1" doc:"The ID of the current resource"`
	Body struct {
		Schedule         string `json:"schedule" validate:"required"`
		Command          string `json:"command" validate:"required"`
		Name             string `json:"name"`
		RunLogsKeepCount int32  `json:"run_logs_keep_count" validate:"gte=0,lte=960"`
		Wrap             bool   `json:"wrap"`
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

type FixUnmanagedCronJobsResponse struct {
	Body struct {
		Fixed uint32 `json:"fixed" doc:"Number of unmanaged cron jobs fixed"`
	}
}

type RunJobRequest struct {
	ID    uint   `path:"id" minimum:"1" doc:"The ID of the current resource"`
	JobID uint64 `path:"job_id" minimum:"1" doc:"The ID of the cron job to run"`
}

type RunJobResponse struct {
	Body *proto.RunCronJobResponse
}
type CronJobInfoRequest struct {
	ID    uint   `path:"id" minimum:"1" doc:"The ID of the current resource"`
	JobID uint64 `path:"job_id" minimum:"1" doc:"The ID of the cron job to run"`
}

type CronJobInfoResponse struct {
	Body *proto.CronInfoResponse
}

type UpdateCronJobRequest struct {
	ID    uint   `path:"id" minimum:"1" doc:"The ID of the current resource"`
	JobID uint64 `path:"job_id" minimum:"1" doc:"The ID of the cron job to update"`
	Body  struct {
		Schedule         string `json:"schedule" validate:"required"`
		Command          string `json:"command" validate:"required"`
		Name             string `json:"name"`
		RunLogsKeepCount int32  `json:"run_logs_keep_count" validate:"gte=0,lte=960"`
		Wrap             bool   `json:"wrap"`
	}
}

type ListCronJobRunsRequest struct {
	ID     uint   `path:"id" minimum:"1" doc:"The ID of the current resource"`
	JobID  uint64 `path:"job_id" minimum:"1" doc:"The ID of the cron job to list runs for"`
	Limit  uint   `query:"limit" default:"10" validate:"lte=100"`
	Offset uint   `query:"offset" default:"0"`
}

type ListCronJobRunsResponse struct {
	Body *proto.ListCronJobRunsResponse
}
