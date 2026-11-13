package resources

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"hazyio.com/servergram/shared/gen/proto"
)

type ResourcesServer struct {
	proto.UnimplementedResourcesServiceServer
}

func (s *ResourcesServer) Get(req *proto.GetServerResourcesRequest, stream proto.ResourcesService_GetServer) error {
	return getResourcesUsage(req, stream)

}
func (s *ResourcesServer) Info(ctx context.Context, _ *emptypb.Empty) (*proto.GetServerInfoResponse, error) {
	return getInfo(ctx)
}
func (s *ResourcesServer) TopProcesses(req *proto.TopProcessRequest, stream proto.ResourcesService_TopProcessesServer) error {
	return topProcesses(req, stream)
}

func (s *ResourcesServer) WatchProcesses(req *proto.WatchProcessesRequest, stream proto.ResourcesService_WatchProcessesServer) error {
	return watchProcesses(req, stream)
}

func (s *ResourcesServer) StopProcess(_ context.Context, req *proto.StopProcessRequest) (*proto.StopProcessResponse, error) {
	return stopProcess(req)
}
func (s *ResourcesServer) ForceStopProcess(_ context.Context, req *proto.ForceStopProcessRequest) (*proto.ForceStopProcessResponse, error) {
	return forceStopProcess(req)
}

func (s *ResourcesServer) SignalProcess(_ context.Context, req *proto.SignalProcessRequest) (*proto.SignalProcessResponse, error) {
	return sendSignal(req)
}