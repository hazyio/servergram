package files

import (
	"context"
	"os"

	"gorm.io/gorm"
	"hazyio.com/servergram/shared/gen/proto"
)

type FilesServer struct {
	database *gorm.DB
	proto.UnimplementedFilesServiceServer
}

func NewFilesServer(database *gorm.DB) *FilesServer {
	return &FilesServer{
		database: database,
	}
}

func (s *FilesServer) ListFiles(req *proto.ListFilesRequest, stream proto.FilesService_ListFilesServer) error {
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
func (s *FilesServer) WatchDirectory(req *proto.WatchDirectoryRequest, stream proto.FilesService_WatchDirectoryServer) error {
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
func (s *FilesServer) PathProperties(_ context.Context, req *proto.PathPropertiesRequest) (*proto.PathPropertiesResponse, error) {
	return handlePathProperties(req)
}

func (s *FilesServer) ReadFile(req *proto.ReadFileRequest, stream proto.FilesService_ReadFileServer) error {
	return handleReadFile(req, stream)
}

func (s *FilesServer) WriteFile(stream proto.FilesService_WriteFileServer) error {
	return handleWriteFile(stream)
}
