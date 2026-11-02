package files

import (
	"context"
	"hazyio.com/servergram/shared/gen/proto"
)

type FilesServer struct {
	proto.UnimplementedFilesServer
}

func (s *FilesServer) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	return &proto.ListFilesResponse{}, nil
}
