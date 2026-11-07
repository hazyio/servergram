package routes

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/routes/models"
)

func registerGroupRoutes(api *huma.Group, database *gorm.DB, allowedOrigins []string, polls *internal.Polls) {
	server := internal.NewGroup("/group", api, database, []string{"Groups"}, nil, allowedOrigins, polls)
	server.Post("/create", http.StatusCreated, createNewGroup)
	server.Get("/list", http.StatusOK, listGroups)
	server.Delete("/delete/{id}", http.StatusOK, deleteGroup)

}

func createNewGroup(data internal.RouteData[models.CreateGroupRequest]) (*EmptyRR, error) {
	input := data.Input
	name := input.Body.Name
	database := data.Database

	find := database.Where("name = ?", name).First(&db.Group{})
	if find.Error == nil {
		return nil, internal.NewHumaError(http.StatusConflict, "Group with that name already exists")
	}
	if find.Error != nil && find.Error != gorm.ErrRecordNotFound {
		return nil, internal.InternalServerError
	}

	newGroup := db.Group{
		Name: name,
	}
	if err := database.Create(&newGroup).Error; err != nil {
		slog.Error("An error occurred while saving group", "error", err)
		return nil, internal.InternalServerError
	}
	db.AddLog(db.LogActionCreate, "Group \""+name+"\" created", data.Database)

	return &EmptyRR{}, nil
}

func listGroups(data internal.RouteData[EmptyRR]) (*models.GroupListResponse, error) {
	database := data.Database

	response := &models.GroupListResponse{}
	response.Body.Items = make([]struct {
		Id                uint      `json:"id"`
		Name              string    `json:"name"`
		CreatedAt         time.Time `json:"created_at"`
		ResourcesAttached uint      `json:"resources_attached"`
	}, 0)

	database.Model(&db.Group{}).
		Select("groups.id, groups.name, groups.created_at, COUNT(servers.id) AS resources_attached").
		Joins("LEFT JOIN servers ON servers.group_id = groups.id").
		Group("groups.id").
		Scan(&response.Body.Items)

	return response, nil
}

func deleteGroup(data internal.RouteData[PathIdRequest]) (*EmptyRR, error) {
	database := data.Database
	id := data.Input.ID

	var group db.Group
	if err := database.First(&group, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, internal.NewHumaError(http.StatusNotFound, "Group not found")
		}
		slog.Error("An error occurred while fetching group", "error", err)
		return nil, internal.InternalServerError
	}
	// Check if there are any servers attached to the group
	var count int64
	if err := database.Model(&db.Server{}).Where("group_id = ?", id).Count(&count).Error; err != nil {
		slog.Error("An error occurred while checking attached servers", "error", err)
		return nil, internal.InternalServerError
	}
	if count > 0 {
		return nil, internal.NewHumaError(http.StatusConflict, "Cannot delete group with attached servers")
	}
	if err := database.Delete(&group).Error; err != nil {
		slog.Error("An error occurred while deleting group", "error", err)
		return nil, internal.InternalServerError
	}
	db.AddLog(db.LogActionDelete, "Group \""+group.Name+"\" deleted", data.Database)

	return &EmptyRR{}, nil
}
