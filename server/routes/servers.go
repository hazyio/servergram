package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/routes/models"
)

func registerServerRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/servers", api, database, []string{"Servers"})
	server.Post("/add", http.StatusCreated, addServer)
	server.Get("/list", http.StatusOK, listServers)
}

func addServer(data internal.NormalRouteData[models.AddServerRequest]) (*models.AddServerResponse, error) {
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

func listServers(data internal.NormalRouteData[EmptyRR]) (*models.ListServersResponse, error) {
	database := data.Database
	var total int64
	if err := database.Model(&db.Server{}).Count(&total).Error; err != nil {
		slog.Error("count servers failed", "error", err)
		return nil, internal.InternalServerError
	}
	response := &models.ListServersResponse{}
	response.Body.Items = make([]models.ServerItem, total)

	database.Model(&db.Server{}).
		Select("servers.id, servers.name, servers.host, servers.created_at").
		Joins("LEFT JOIN groups ON servers.group_id = groups.id").
		Group("servers.id").
		Scan(&response.Body.Items)
	response.Body.Total = total

	return response, nil
}
