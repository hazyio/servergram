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

	server.SSE("/list-in-server", http.StatusOK, listFilesInServer)

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

func listFilesInServer(data internal.RouteData[models.FileListRequest], send sse.Sender) *EmptyRR {
	sshManager := data.SSHManager
	ctx := data.Context
	if sshManager == nil {
		send.Data(internal.NewSSEError(http.StatusInternalServerError, "SSH manager not initialized"))
		return &EmptyRR{}
	}
	resource, ctx, cancel, err := getFilesConnected(ctx, sshManager, data.Input.Body.Id)
	if err != nil {
		slog.Error("Failed to get files connected", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))
		return &EmptyRR{}
	}
	defer cancel()
	stream, err := resource.ListFiles(ctx, &proto.ListFilesRequest{
		Dir:           data.Input.Body.Dir,
		IncludeHidden: data.Input.Body.IncludeHidden,
	})
	if err != nil {

		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, "Internal server error"))
		return &EmptyRR{}
	}
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			send.Data(internal.NewSSEInternalError(err.Error()))

			break // server finished the stream cleanly
		}
		if err != nil {
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.NotFound:
				send.Data(internal.NewSSEError(http.StatusNotFound, err.Error()))
			case codes.PermissionDenied: // 403
				send.Data(internal.NewSSEError(http.StatusForbidden, err.Error()))

			default:
				send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))

			}
			return &EmptyRR{}
		}
		send.Data(chunk.Entries)
		if chunk.Done {
			break
		}
	}
	return &EmptyRR{}
}
