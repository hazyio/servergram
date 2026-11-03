package resources

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"hazyio.com/servergram/shared/gen/proto"
)

type ResourcesServer struct {
	proto.UnimplementedServerResourcesServer
}

func (s *ResourcesServer) Get(ctx context.Context, _ *emptypb.Empty) (*proto.GetServerResourcesResponse, error) {
	memUsage, err := memoryUsage()
	if err != nil {
		return nil, err
	}
	cpuUsage, err := cpuUsage(ctx)
	if err != nil {
		return nil, err
	}
	diskUsage, err := diskUsage(ctx)
	if err != nil {
		return nil, err
	}
	swapUsage, err := swapUsage(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetServerResourcesResponse{
		Cpu:    cpuUsage,
		Memory: &memUsage,
		Disk:   diskUsage,
		Swap:   swapUsage,
	}, nil
}
func (s *ResourcesServer) Info(ctx context.Context, _ *emptypb.Empty) (*proto.GetServerInfoResponse, error) {
	return getInfo(ctx)
}
func (s *ResourcesServer) TopProcesses(ctx context.Context, req *proto.TopProcessRequest) (*proto.TopProcessResponse, error) {
	return topProcesses(ctx, req)
}
