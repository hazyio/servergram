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

}

func createNewGroup(data internal.NormalRouteData[CreateGroupRequest]) (*EmptyResponse, error) {
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

	return &EmptyResponse{}, nil
}
