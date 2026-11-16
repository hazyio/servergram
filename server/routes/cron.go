package routes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
	"hazyio.com/servergram/server/routes/models"
	"hazyio.com/servergram/shared/gen/proto"
)

func registerCronRoutes(api *huma.Group, database *gorm.DB, ssh *sshmanager.SSHManager, allowedOrigins []string, polls *internal.Polls) {
	server := internal.NewGroup("/cron", api, database, []string{"Cron"}, ssh, allowedOrigins, polls)
	server.SSE("/{id}/list", http.StatusOK, listCronJobs)
	server.Post("/{id}/create", http.StatusOK, createCronJob)
	server.Post("/{id}/remove-duplicates", http.StatusOK, removeDuplicateCronJobs)
}
func getCronConnected(ctx context.Context, manager *sshmanager.SSHManager, id uint) (proto.CronServiceClient, context.Context, context.CancelFunc, error) {
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
	resource, ctx, cancel := agentConnection.Agent.CronClient(ctx)

	return resource, ctx, cancel, nil
}
func listCronJobs(data internal.RouteData[PathIdRequest], send sse.Sender) *proto.ListCronJobsResponse {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		slog.Error("Failed to get files connected", "error", err)
		_ = send.Data(internal.NewSSEError(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer cancel()
	stream, err := resource.ListCronJobs(ctx, &proto.ListCronJobsRequest{})
	if err != nil {
		slog.Error("Failed to list cron jobs", "error", err)
		internal.SendSSEError(send, err)
		return nil
	}
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Listed cron jobs in server \"%d\"", data.Input.ID), data.Database)

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Stream ended without a done chunk; tell the client to stop reconnecting.
			_ = send.Data(&proto.ListCronJobsResponse{Done: true})
			return nil

		}
		if err != nil {
			slog.Error("Failed to receive cron job chunk", "error", err)
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
func createCronJob(data internal.RouteData[models.CreateCronJobRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.CreateCronJob(ctx, &proto.CreateCronJobRequest{
		Schedule:          data.Input.Body.Schedule,
		Command:           data.Input.Body.Command,
		Name:              data.Input.Body.Name,
		MetaKeepCount:     data.Input.Body.MetaKeepCount,
		RunLogsKeepCount:  data.Input.Body.RunLogsKeepCount,
		ErrorLogsKeepDays: data.Input.Body.ErrorLogsKeepDays,
		Wrap:              data.Input.Body.Wrap,
	})
	if err != nil {
		slog.Error("Failed to create cron job", "error", err)
		return nil, err
	}
	db.AddLog(db.LogActionCreate, fmt.Sprintf("Created cron job in server \"%d\"", data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}

func removeDuplicateCronJobs(data internal.RouteData[PathIdRequest]) (*models.RemoveDuplicateCronJobsResponse, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()

	result, err := resource.RemoveDuplicateCronJobs(ctx, &proto.RemoveDuplicateCronJobsRequest{})
	if err != nil {
		slog.Error("Failed to remove duplicate cron jobs", "error", err)
		return nil, err
	}

	response := &models.RemoveDuplicateCronJobsResponse{}
	response.Body.Removed = result.Removed
	db.AddLog(db.LogActionDelete, fmt.Sprintf("Removed %d duplicate cron jobs in server \"%d\"", result.Removed, data.Input.ID), data.Database)
	return response, nil
}
