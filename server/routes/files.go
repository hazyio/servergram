package routes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

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
	server.Get("/get-path-properties-in-server/{id}", http.StatusOK, getPathPropertiesInServer)
	server.SSE("/read-file-in-server/{id}", http.StatusOK, readFileInServer)
	server.Download("/download-file-in-server/{id}", downloadFileInServer, http.StatusBadGateway)
	server.Upload("/upload-file-to-server/{id}", uploadFileInServer, http.StatusBadGateway)

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

func getPathPropertiesInServer(data internal.RouteData[models.PathPropertiesRequest]) (*models.PathPropertiesResponse, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	resp, err := resource.PathProperties(ctx, &proto.PathPropertiesRequest{
		Path: data.Input.Path,
		Name: data.Input.Name,
	})
	if err != nil {
		slog.Error("Failed to get path properties from server", "error", err)
		return nil, internal.FromGrpcError(err)
	}
	result := &models.PathPropertiesResponse{
		Body: struct {
			Enity *proto.FileEntry
		}{
			Enity: resp.Entry,
		},
	}
	return result, nil
}
func readFileInServer(data internal.RouteData[models.ReadFileRequest], send sse.Sender) *proto.ReadFileResponse {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		slog.Error("Failed to get files connected", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer cancel()

	stream, err := resource.ReadFile(ctx, &proto.ReadFileRequest{
		Path:         data.Input.Path,
		Name:         data.Input.Name,
		MaxBytesSize: data.Input.MaxBytesSize,
		IsDownload:   false,
	})
	if err != nil {
		slog.Error("Failed to read file", "error", err)
		internal.SendSSEError(send, err)
		return nil
	}

	db.AddLog(db.LogActionUpdate, fmt.Sprintf("read file \"%s\" in server \"%d\"", data.Input.Name, data.Input.ID), data.Database)

	for {
		select {
		case <-stream.Context().Done():
			slog.Debug("Stream context cancelled")
			return nil

		default:
			chunk, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// Stream ended without a done chunk; tell the client to stop reconnecting.
				_ = send.Data(&proto.ReadFileResponse{Done: true})
				return nil
			}
			if err != nil {
				slog.Error("Failed to receive file chunk", "error", err)

				internal.SendSSEError(send, err)
				return nil
			}
			slog.Debug("Read file chunk", "chunk_size", len(chunk.Data))

			if err := send.Data(chunk); err != nil {
				slog.Error("Failed to send file chunk", "error", err)
				return nil
			}
			if chunk.Done {
				return nil
			}
		}
	}
}
func downloadFileInServer(data internal.RouteData[models.DownloadFileRequest]) (*huma.StreamResponse, error) {
	getMetaResource, getMetaCtx, getMetaCancel, getMetaErr := getFilesConnected(context.Background(), data.SSHManager, data.Input.ID)
	if getMetaErr != nil {
		slog.Error("Failed to get files connected", "error", getMetaErr)
		return nil, getMetaErr
	}
	defer getMetaCancel()

	meta, err := getMetaResource.PathProperties(getMetaCtx, &proto.PathPropertiesRequest{
		Path: data.Input.Path,
		Name: data.Input.Name,
	})
	if err != nil {
		slog.Error("Failed to get file properties", "error", err)
		return nil, internal.Error502BadGateway
	}

	db.AddLog(db.LogActionUpdate, fmt.Sprintf("download file \"%s\" in server \"%d\"", data.Input.Name, data.Input.ID), data.Database)

	return &huma.StreamResponse{
		Body: func(hctx huma.Context) {

			hctx.SetHeader("Content-Type", "application/octet-stream")
			hctx.SetHeader("Content-Disposition",
				mime.FormatMediaType("attachment", map[string]string{"filename": meta.Entry.Name}))
			if meta.Entry.Size > 0 {
				hctx.SetHeader("Content-Length", strconv.FormatInt(meta.Entry.Size, 10))
			}

			w := hctx.BodyWriter()
			// the previous connect will be cancel, here we use the route context
			resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
			if err != nil {
				slog.Error("Failed to get files connected for streaming", "error", err)
				panic(http.ErrAbortHandler) // headers already sent, so it stops
			}
			defer cancel()
			stream, err := resource.ReadFile(ctx, &proto.ReadFileRequest{
				Path:       data.Input.Path,
				Name:       data.Input.Name,
				IsDownload: true,
			})
			if err != nil {
				slog.Error("Failed to start download stream", "error", err)
				panic(http.ErrAbortHandler) // headers already sent, so it stops
			}
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					return
				}

				if err != nil {
					slog.Error("Failed to receive download chunk", "error", err)
					panic(http.ErrAbortHandler) // headers already sent
				}
				if _, err := w.Write(msg.Data); err != nil {
					return // client disconnected
				}
			}
		},
	}, nil
}

func uploadFileInServer(data internal.RouteData[models.DownloadFileRequest], reader io.Reader) (*EmptyRR, error) {
	resource, ctx, cancel, err := getFilesConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		slog.Error("Failed to get files connected for streaming", "error", err)
		return nil, err
	}
	defer cancel()
	stream, err := resource.WriteFile(ctx)
	if err != nil {
		slog.Error("Failed to start upload stream", "error", err)
		return nil, internal.InternalServerError
	}
	stream.Send(&proto.WriteFileRequest{
		Payload: &proto.WriteFileRequest_Meta{
			Meta: &proto.WriteFileRequest_WriteFileMeta{Name: data.Input.Name, Path: data.Input.Path},
		},
	})

	return &EmptyRR{}, nil
}
