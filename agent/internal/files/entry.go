package files

import (
	"context"
	pb "hazyio.com/servergram/shared/gen/proto"
)

type FilesServer struct {
	pb.UnimplementedFilesServer
}

func (s *FilesServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// real implementation here
	return &pb.ListFilesResponse{}, nil
}
