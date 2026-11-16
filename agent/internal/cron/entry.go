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
	return handleListCronJobs(s.database, req, stream)
}

func (s *CronServer) CreateCronJob(ctx context.Context, req *proto.CreateCronJobRequest) (*proto.CreateCronJobResponse, error) {
	return handleCreateCronJob(ctx, s.database, req)
}
