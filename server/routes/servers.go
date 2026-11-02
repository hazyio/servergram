package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
	"hazyio.com/servergram/server/routes/models"
)

func registerServerRoutes(api *huma.Group, database *gorm.DB, sshManager *sshmanager.SSHManager) {
	server := internal.NewGroup("/server", api, database, []string{"Servers"}, sshManager)
	server.Post("/add", http.StatusCreated, addServer, http.StatusConflict, http.StatusInternalServerError)
	server.Get("/list", http.StatusOK, listServers, http.StatusInternalServerError)
	server.Delete("/delete/{id}", http.StatusOK, deleteServer, http.StatusInternalServerError, http.StatusInternalServerError, http.StatusNotFound)
	server.Put("/update-finger-print/{id}", http.StatusNoContent, updateServerFingerPrint, http.StatusNotFound, http.StatusInternalServerError, http.StatusNotFound)
	server.Get("/watch-resources/{id}", http.StatusOK, watchServerResource, http.StatusInternalServerError, http.StatusNotFound, http.StatusForbidden)

}

func addServer(data internal.RouteData[models.AddServerRequest]) (*models.AddServerResponse, error) {
	input := data.Input.Body
	database := data.Database
	_, err := db.FindSSHKeyByID(input.SSHKeyID, database)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, internal.NewHumaError(http.StatusNotFound, "SSH key does not exist")
	}
	if err != nil {
		return nil, internal.InternalServerError
	}
	_, err = db.FindGroupByID(input.GroupID, database)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, internal.NewHumaError(http.StatusNotFound, "Group does not exist")
	}
	if err != nil {
		return nil, internal.InternalServerError
	}

	if _, err := db.FindServerByName(input.Name, database); err == nil {
		return nil, internal.NewHumaError(http.StatusConflict, "Server with this name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, internal.InternalServerError
	}
	newServer := &db.Server{
		Name:               input.Name,
		Host:               input.Host,
		Port:               input.Port,
		Username:           input.UserName,
		OS:                 db.UnknownOS,
		Arch:               db.UnknownArch,
		GroupID:            input.GroupID,
		SSHKeyID:           input.SSHKeyID,
		HostKeyFingerprint: "",
		Trusted:            false,
	}
	if err := database.Create(newServer).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, internal.NewHumaError(http.StatusConflict, "Server with this name, host and username already exists")
		}
		return nil, internal.InternalServerError
	}
	db.AddLog(db.Add, "Added server"+newServer.Name+"", database)
	return &models.AddServerResponse{
		Body: struct {
			ID uint `json:"id" doc:"The ID of the newly created server"`
		}{
			ID: newServer.ID,
		},
	}, nil
}

func listServers(data internal.RouteData[EmptyRR]) (*models.ListServersResponse, error) {
	database := data.Database

	response := &models.ListServersResponse{}
	response.Body.Items = make([]models.ServerItem, 0)

	database.Model(&db.Server{}).
		Select("servers.id, servers.name, servers.host, servers.created_at, servers.trusted").Order("servers.id").
		Scan(&response.Body.Items)

	return response, nil
}

func deleteServer(data internal.RouteData[PathIdRequest]) (*EmptyRR, error) {
	database := data.Database
	id := data.Input.ID
	server, err := db.FindServerByID(id, database)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, internal.NewHumaError(http.StatusNotFound, "Server does not exist")
	}
	if err != nil {
		return nil, internal.InternalServerError
	}
	if err := database.Delete(server).Error; err != nil {
		return nil, internal.InternalServerError
	}
	db.AddLog(db.Delete, "Deleted server "+server.Name+"", database)
	return &EmptyRR{}, nil
}

func updateServerFingerPrint(data internal.RouteData[models.UpdateServerFingerPrintRequest]) (*EmptyRR, error) {
	database := data.Database
	id := data.Input.ID
	var server db.Server
	if err := database.First(&server, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, internal.NewHumaError(http.StatusNotFound, "Server not found")
		}
		return nil, internal.InternalServerError
	}

	newFingerPrint := data.Input.Body.Fingerprint
	server.HostKeyFingerprint = newFingerPrint
	server.Trusted = true

	if err := database.Save(&server).Error; err != nil {
		return nil, internal.InternalServerError
	}
	db.AddLog(db.Update, "Updated host fingerprint for server "+server.Name+"", database)
	return &EmptyRR{}, nil
}

func watchServerResource(data internal.RouteData[models.WatchServerResourceRequest]) (*models.WatchServerResourceResponse, error) {
	// database := data.Database
	id := data.Input.ID
	ctx := data.Context
	sshmanager := data.SSHManager
	if sshmanager == nil {
		return nil, internal.NewHumaError(http.StatusInternalServerError, "SSH Manager not initialized")
	}
	if !sshmanager.IsConnected(id) {
		return nil, internal.NewHumaError(http.StatusNotFound, "Server is not connected")
	}
	agentConnection, err := sshmanager.GetConnected(id)
	if err != nil {
		return nil, internal.NewHumaError(http.StatusInternalServerError, "Failed to get connected agent, it is probably not connected")
	}
	resource, ctx, cancel := agentConnection.Agent.ServerResourcesClient(ctx)
	defer cancel()
	res, err := resource.Get(ctx, &emptypb.Empty{})
	if err != nil {
		slog.Error("Failed to get server resources", "error", err)
		return nil, internal.NewHumaError(http.StatusInternalServerError, "Failed to get server resources")
	}
	slog.Info("Got server resources,", "resources", res)
	response := &models.WatchServerResourceResponse{}
	response.Body = res

	return response, nil
}
