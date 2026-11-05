package resources

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"hazyio.com/servergram/shared/gen/proto"
)

type ResourcesServer struct {
	proto.UnimplementedServerResourcesServer
}

func (s *ResourcesServer) Get(req *proto.GetServerResourcesRequest, stream proto.ServerResources_GetServer) error {
	return getResourcesUsage(req, stream)

}
func (s *ResourcesServer) Info(ctx context.Context, _ *emptypb.Empty) (*proto.GetServerInfoResponse, error) {
	return getInfo(ctx)
}
func (s *ResourcesServer) TopProcesses(req *proto.TopProcessRequest, stream proto.ServerResources_TopProcessesServer) error {
	return topProcesses(req, stream)
}
	