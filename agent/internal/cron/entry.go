package cron

import (
	"context"

	"gorm.io/gorm"

	"hazyio.com/servergram/shared/gen/proto"
)

type CronServer struct {
	database *gorm.DB
	proto.UnimplementedCronServiceServer
}

func NewCronServer(database *gorm.DB) *CronServer {
	return &CronServer{
		database: database,
	}
}

func (s *CronServer) ListCronJobs(req *proto.ListCronJobsRequest, stream proto.CronService_ListCronJobsServer) error {
	return handleListCronJobs(s.database, stream)
}

func (s *CronServer) CreateCronJob(ctx context.Context, req *proto.CreateCronJobRequest) (*proto.CreateCronJobResponse, error) {
	return handleCreateCronJob(ctx, s.database, req)
}

func (s *CronServer) RemoveDuplicateCronJobs(ctx context.Context, req *proto.RemoveDuplicateCronJobsRequest) (*proto.RemoveDuplicateCronJobsResponse, error) {
	return handleRemoveDuplicateCronJobs(ctx)
}

func (s *CronServer) FixUnmanagedCronJobs(ctx context.Context, req *proto.FixUnmanagedCronJobsRequest) (*proto.FixUnmanagedCronJobsResponse, error) {
	return handleFixUnmanagedCronJobs(ctx, s.database)
}

func (s *CronServer) RunCronJob(ctx context.Context, req *proto.RunCronJobRequest) (*proto.RunCronJobResponse, error) {

	return handleRunCronJob(ctx, s.database, req)
}
func (s *CronServer) StopCronJob(ctx context.Context, req *proto.StopCronJobRequest) (*proto.StopCronJobResponse, error) {
	return handleStopCronJob(ctx, s.database, req)
}

func (s *CronServer) CronInfo(ctx context.Context, req *proto.CronInfoRequest) (*proto.CronInfoResponse, error) {
	return handleInfo(ctx, s.database, req)
}

func (s *CronServer) UpdateCronJob(ctx context.Context, req *proto.UpdateCronJobRequest) (*proto.UpdateCronJobResponse, error) {
	return handleUpdateCronJob(ctx, s.database, req)
}
