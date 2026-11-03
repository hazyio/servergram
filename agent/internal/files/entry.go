package files

import (
	"hazyio.com/servergram/shared/gen/proto"
)

type FilesServer struct {
	proto.UnimplementedFilesServer
}

func (s *FilesServer) ListFiles(req *proto.ListFilesRequest, stream proto.Files_ListFilesServer) error {
	
	return HandleListFiles(req, stream)
}
