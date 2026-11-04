package files

import (
	"os"

	"hazyio.com/servergram/shared/gen/proto"
)

type FilesServer struct {
	proto.UnimplementedFilesServer
}

func (s *FilesServer) ListFiles(req *proto.ListFilesRequest, stream proto.Files_ListFilesServer) error {
	dir := req.Dir
	if dir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			dir = homeDir
		} else {
			dir = "."
		}
	}
	req.Dir = dir
	return HandleListFiles(req, stream)
}
