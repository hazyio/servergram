package routes

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
)

func registerGroupRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/group", api, database)
	server.Post("/create", http.StatusCreated, createNewGroup)
	server.Get("/list", http.StatusOK, listGroups)

}

func createNewGroup(data internal.NormalRouteData[CreateGroupRequest]) (*EmptyRR, error) {
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

	return &EmptyRR{}, nil
}

func listGroups(data internal.NormalRouteData[EmptyRR]) (*GroupListResponse, error) {
	database := data.Database
	var groups []db.Group
	if err := database.Find(&groups).Error; err != nil {
		slog.Error("An error occurred while fetching groups", "error", err)
		return nil, internal.InternalServerError
	}
	response := &GroupListResponse{}
	response.Body.Items = make([]struct {
		Id                uint   `json:"id"`
		Name              string `json:"name"`
		CreatedAt         string `json:"created_at"`
		ResourcesAttached uint   `json:"resources_attached"`
	}, len(groups))
	for i, group := range groups {
		response.Body.Items[i].Id = group.ID
		response.Body.Items[i].Name = group.Name
		response.Body.Items[i].CreatedAt = group.CreatedAt.String()
		response.Body.Items[i].ResourcesAttached = 0
	}
	response.Body.Total = len(groups)

	return response, nil
}
