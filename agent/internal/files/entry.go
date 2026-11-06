package files

import (
	"context"
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
		if err != nil {
			return err
		}
		dir = homeDir
	}
	req.Dir = dir
	return handleListFiles(req, stream)
}

func (s *FilesServer) CreateDirectory(_ context.Context, req *proto.CreateDirectoryRequest) (*proto.SuccessResponse, error) {
	return handleCreateDirectory(req)
}
func (s *FilesServer) RemoveDirectory(_ context.Context, req *proto.RemoveDirectoryRequest) (*proto.SuccessResponse, error) {
	return handleRemoveDirectory(req)
}
func (s *FilesServer) WatchDirectory(req *proto.WatchDirectoryRequest, stream proto.Files_WatchDirectoryServer) error {
	return handleWatchDirectory(req, stream)
}
func (s *FilesServer) CreateFile(_ context.Context, req *proto.CreateFileRequest) (*proto.SuccessResponse, error) {
	return handleCreateFile(req)
}
func (s *FilesServer) RemoveFile(_ context.Context, req *proto.RemoveFileRequest) (*proto.SuccessResponse, error) {
	return handleDeleteFile(req)
}
func (s *FilesServer) RenamePath(_ context.Context, req *proto.RenamePathRequest) (*proto.SuccessResponse, error) {
	return handleRenamePath(req)
}
