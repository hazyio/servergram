package resources

import (
	"context"
	"hazyio.com/servergram/shared/gen/proto"
)

type ResourcesServer struct {
	proto.UnimplementedServerResourcesServer
}

func (s *ResourcesServer) Get(ctx context.Context, req *proto.GetServerResourcesRequest) (*proto.GetServerResourcesResponse, error) {
	return &proto.GetServerResourcesResponse{}, nil
}
