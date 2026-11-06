package routes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
	"hazyio.com/servergram/server/routes/models"
	"hazyio.com/servergram/shared/gen/proto"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

func registerFilesRoutes(api *huma.Group, database *gorm.DB, ssh *sshmanager.SSHManager) {
	server := internal.NewGroup("/files", api, database, []string{"Files"}, ssh)

	server.SSE("/list-in-server/{id}", http.StatusOK, listFilesInServer)
	server.SSE("/watch-dir-in-server/{id}", http.StatusOK, watchDirectoryInServer)
	server.Put("/create-dir-in-server/{id}", http.StatusOK, createDirectoryInServer)
	server.Delete("/remove-dir-in-server/{id}", http.StatusOK, removeDirectoryInServer)
	server.Put("/create-file-in-server/{id}", http.StatusOK, createFileInServer)
	server.Delete("/remove-file-in-server/{id}", http.StatusOK, removeFileInServer)
	server.Patch("/rename-path-in-server/{id}", http.StatusOK, renamePathInServer)
	
}
func getFilesConnected(ctx context.Context, manager *sshmanager.SSHManager, id uint) (proto.FilesClient, context.Context, context.CancelFunc, error) {
	if manager == nil {
		return nil, nil, nil, internal.NewHumaError(http.StatusInternalServerError, "SSH Manager not initialized")
	}
	if !manager.IsConnected(id) {
		return nil, nil, nil, internal.NewHumaError(http.StatusNotFound, "Server is not connected")
	}
	agentConnection, err := manager.GetConnected(id)
	if err != nil {
		return nil, nil, nil, internal.NewHumaError(http.StatusInternalServerError, "Failed to get connected agent, it is probably not connected")
	}
	resource, ctx, cancel := agentConnection.Agent.FilesClient(ctx)

	return resource, ctx, cancel, nil
}

func listFilesInServer(data internal.RouteData[models.FileListRequest], send sse.Sender) *proto.ListFilesChunkResponse {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		slog.Error("Failed to get files connected", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer cancel()
	stream, err := resource.ListFiles(ctx, &proto.ListFilesRequest{
		Dir:           data.Input.Dir,
		IncludeHidden: data.Input.IncludeHidden,
	})
	if err != nil {
		slog.Error("Failed to list files", "error", err)
		internal.SendSSEError(send, err)
		return nil
	}
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Listed files in server \"%d\", dir \"%s\"", data.Input.ID, data.Input.Dir), data.Database)

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Stream ended without a done chunk; tell the client to stop reconnecting.
			_ = send.Data(&proto.ListFilesChunkResponse{Done: true})
			return nil

		}
		if err != nil {
			slog.Error("Failed to receive file chunk", "error", err)
			internal.SendSSEError(send, err)
			return nil
		}
		if err := send.Data(chunk); err != nil {
			return nil
		}
		if chunk.Done {
			return nil
		}
	}
}

func watchDirectoryInServer(data internal.RouteData[models.WatchDirectoryRequest], send sse.Sender) *proto.WatchDirectoryResponse {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		slog.Error("Failed to get files connected", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer cancel()
	stream, err := resource.WatchDirectory(ctx, &proto.WatchDirectoryRequest{
		Path:          data.Input.Dir,
		IncludeHidden: data.Input.IncludeHidden,
	})
	if err != nil {
		internal.SendSSEError(send, err)
		return nil
	}

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			internal.SendSSEError(send, err)
			return nil
		}
		if err := send.Data(chunk); err != nil {
			return nil
		}

	}
}

func createDirectoryInServer(data internal.RouteData[models.CreateDirectoryRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.CreateDirectory(ctx, &proto.CreateDirectoryRequest{
		Path: data.Input.Body.Path,
		Name: data.Input.Body.Name,
	})
	if err != nil {
		slog.Error("Failed to create directory in server", "error", err)
		return nil, internal.FromGrpcError(err)

	}
	db.AddLog(db.LogActionCreate, fmt.Sprintf("Created directory \"%s\" in server \"%d\"", data.Input.Body.Name, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
func removeDirectoryInServer(data internal.RouteData[models.RemoveDirectoryRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.RemoveDirectory(ctx, &proto.RemoveDirectoryRequest{
		Path: data.Input.Body.Path,
	})
	if err != nil {
		slog.Error("Failed to remove directory from server", "error", err)
		return nil, internal.FromGrpcError(err)

	}
	db.AddLog(db.LogActionDelete, fmt.Sprintf("Deleted directory \"%s\" from server \"%d\"", data.Input.Body.Path, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
func createFileInServer(data internal.RouteData[models.CreateFileRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.CreateFile(ctx, &proto.CreateFileRequest{
		Path: data.Input.Body.Path,
		Name: data.Input.Body.Name,
	})
	if err != nil {
		slog.Error("Failed to create file on server", "error", err)
		return nil, internal.FromGrpcError(err)

	}
	db.AddLog(db.LogActionCreate, fmt.Sprintf("Created file \"%s\" in server \"%d\"", data.Input.Body.Name, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
func removeFileInServer(data internal.RouteData[models.DeleteFileRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.RemoveFile(ctx, &proto.RemoveFileRequest{
		Path: data.Input.Body.Path,
		Name: data.Input.Body.Name,
	})
	if err != nil {

		slog.Error("Failed to remove file from server", "error", err)
		return nil, internal.FromGrpcError(err)

	}

	db.AddLog(db.LogActionDelete, fmt.Sprintf("Deleted file \"%s\" from server \"%d\"", data.Input.Body.Name, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
func renamePathInServer(data internal.RouteData[models.RenamePathRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.RenamePath(ctx, &proto.RenamePathRequest{
		Path:    data.Input.Body.Path,
		OldName: data.Input.Body.OldName,
		NewName: data.Input.Body.NewName,
	})

	if err != nil {
		slog.Error("Failed to rename path on server", "error", err)
		return nil, internal.FromGrpcError(err)
	}

	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Renamed path \"%s\" to \"%s\" in server \"%d\"", data.Input.Body.OldName, data.Input.Body.NewName, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
