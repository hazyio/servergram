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
	server.Post("/{id}/fix-unmanaged", http.StatusOK, fixUnmanagedCronJobs)
	server.Get("/{id}/run/{job_id}", http.StatusOK, runCronJob)
	server.Get("/{id}/info/{job_id}", http.StatusOK, getCronJobInfo)
	server.Post("/{id}/update/{job_id}", http.StatusOK, updateCronJob)
	server.Delete("/{id}/delete/{job_id}", http.StatusOK, deleteCronJob)
	server.Put("/{id}/duplicate/{job_id}", http.StatusOK, duplicateCronJob)
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

func fixUnmanagedCronJobs(data internal.RouteData[PathIdRequest]) (*models.FixUnmanagedCronJobsResponse, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()

	result, err := resource.FixUnmanagedCronJobs(ctx, &proto.FixUnmanagedCronJobsRequest{})
	if err != nil {
		slog.Error("Failed to fix unmanaged cron jobs", "error", err)
		return nil, err
	}
	response := &models.FixUnmanagedCronJobsResponse{}
	response.Body.Fixed = result.Fixed
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Fixed %d unmanaged cron jobs in server \"%d\"", result.Fixed, data.Input.ID), data.Database)
	return response, nil
}

func runCronJob(data internal.RouteData[models.RunJobRequest]) (*models.RunJobResponse, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	result, err := resource.RunCronJob(ctx, &proto.RunCronJobRequest{
		Id: data.Input.JobID,
	})
	if err != nil {
		slog.Error("Failed to run cron job", "error", err)
		return nil, err
	}
	response := &models.RunJobResponse{}
	response.Body = result
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Ran cron job %d in server \"%d\"", data.Input.JobID, data.Input.ID), data.Database)
	return response, nil
}

func getCronJobInfo(data internal.RouteData[models.CronJobInfoRequest]) (*models.CronJobInfoResponse, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	result, err := resource.CronInfo(ctx, &proto.CronInfoRequest{
		Id: data.Input.JobID,
	})
	if err != nil {
		slog.Error("Failed to get cron job info", "error", err)
		return nil, err
	}
	response := &models.CronJobInfoResponse{}
	response.Body = result
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Retrieved cron job %d info in server \"%d\"", data.Input.JobID, data.Input.ID), data.Database)
	return response, nil
}

func updateCronJob(data internal.RouteData[models.UpdateCronJobRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}

	defer cancel()
	_, err = resource.UpdateCronJob(ctx, &proto.UpdateCronJobRequest{
		Id:                data.Input.JobID,
		Schedule:          data.Input.Body.Schedule,
		Command:           data.Input.Body.Command,
		Name:              data.Input.Body.Name,
		MetaKeepCount:     data.Input.Body.MetaKeepCount,
		RunLogsKeepCount:  data.Input.Body.RunLogsKeepCount,
		ErrorLogsKeepDays: data.Input.Body.ErrorLogsKeepDays})
	if err != nil {
		slog.Error("Failed to update cron job", "error", err)
		return nil, err
	}
	db.AddLog(db.LogActionUpdate, fmt.Sprintf("Updated cron job %d in server \"%d\"", data.Input.JobID, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}

func deleteCronJob(data internal.RouteData[models.CronJobInfoRequest]) (*EmptyRR, error) {
	resource, ctx, cancel, err := getCronConnected(data.Context, data.SSHManager, data.Input.ID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = resource.DeleteCronJob(ctx, &proto.DeleteCronJobRequest{
		Id: data.Input.JobID,
	})
	if err != nil {
		slog.Error("Failed to delete cron job", "error", err)
		return nil, err
	}
	db.AddLog(db.LogActionDelete, fmt.Sprintf("Deleted cron job %d from server \"%d\"", data.Input.JobID, data.Input.ID), data.Database)
	return &EmptyRR{}, nil
}
