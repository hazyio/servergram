package files

import (
	"errors"
	"io/fs"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/shared/gen/proto"
)

func HandleListFiles(req *proto.ListFilesRequest, stream proto.Files_ListFilesServer) error {
	entries, err := os.ReadDir(req.Dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return status.Errorf(codes.NotFound, "directory %q not found", req.Dir)
		}
		if errors.Is(err, fs.ErrPermission) {
			return status.Errorf(codes.PermissionDenied, "cannot read %q", req.Dir)
		}

		return status.Error(codes.Internal, "unknown error")
	}
	batch := make([]*proto.FileEntry, 0, 500)
	for _, e := range entries {
		if e.Name()[0] == '.' && !req.IncludeHidden {
			continue
		}
		info, _ := e.Info()
		batch = append(batch, toEntry(info))
		if len(batch) == 500 {
			stream.Send(&proto.ListFilesChunkResponse{Entries: batch})
			batch = batch[:0]
		}
	}
	stream.Send(&proto.ListFilesChunkResponse{Entries: batch, Done: true})

	return nil
}
