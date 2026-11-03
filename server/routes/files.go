package routes

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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

func sendFilesError(send sse.Sender, err error) {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.NotFound:
		_ = send.Data(internal.NewSSEError(http.StatusNotFound, st.Message()))
	case codes.PermissionDenied:
		_ = send.Data(internal.NewSSEError(http.StatusForbidden, st.Message()))
	default:
		slog.Error("Failed to list files", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, "Internal server error"))
	}
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
		sendFilesError(send, err)
		return nil
	}
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Stream ended without a done chunk; tell the client to stop reconnecting.
			_ = send.Data(&proto.ListFilesChunkResponse{Done: true})
			return nil
		}
		if err != nil {
			sendFilesError(send, err)
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
