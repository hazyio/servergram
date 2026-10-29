package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/routes/models"
)

func registerServerRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/servers", api, database)
	server.Post("/ssh", http.StatusOK, searchSSHKeys)
	server.Post("/group", http.StatusOK, searchGroups)
}

func addServer(data internal.NormalRouteData[models.SearchRequest]) (*models.SearchResponse[models.SSHKeySearchData], error) {
	return nil, internal.InternalServerError
}
