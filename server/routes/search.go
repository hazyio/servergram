package routes

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/routes/models"
)

func registerSearchRoutes(api *huma.Group, database *gorm.DB, allowedOrigins []string) {
	server := internal.NewGroup("/search", api, database, []string{"Search"}, nil, allowedOrigins, nil)
	server.Get("/ssh", http.StatusOK, searchSSHKeys)
	server.Get("/group", http.StatusOK, searchGroups)
}

func searchSSHKeys(data internal.RouteData[models.SearchRequest]) (*models.SearchResponse[models.SSHKeySearchData], error) {
	input := data.Input
	query := input.Query
	database := data.Database

	var results []db.SSHKey

	if err := database.Where("name LIKE ?", "%"+query+"%").Limit(10).Find(&results).Error; err != nil {
		slog.Error("An error occurred while searching SSH keys", "error", err)
		return nil, internal.InternalServerError
	}

	var searchData []models.SSHKeySearchData
	for _, r := range results {
		searchData = append(searchData, models.SSHKeySearchData{
			Id:        r.ID,
			Name:      r.Name,
			Encrypted: r.Encrypted,
		})
	}

	return &models.SearchResponse[models.SSHKeySearchData]{Body: struct {
		Items []models.SSHKeySearchData `json:"items" required:"true" doc:"The search results"`
		Total int                       `json:"total" required:"true" doc:"The total number of search results"`
	}{Items: searchData, Total: len(searchData)}}, nil
}

func searchGroups(data internal.RouteData[models.SearchRequest]) (*models.SearchResponse[models.GroupSearchData], error) {
	input := data.Input
	query := input.Query
	database := data.Database

	var results []db.Group

	if err := database.Where("name LIKE ?", "%"+query+"%").Limit(10).Find(&results).Error; err != nil {
		slog.Error("An error occurred while searching groups", "error", err)
		return nil, internal.InternalServerError
	}

	var searchData []models.GroupSearchData
	for _, r := range results {
		searchData = append(searchData, models.GroupSearchData{
			Id:   r.ID,
			Name: r.Name,
		})
	}

	return &models.SearchResponse[models.GroupSearchData]{Body: struct {
		Items []models.GroupSearchData `json:"items" required:"true" doc:"The search results"`
		Total int                      `json:"total" required:"true" doc:"The total number of search results"`
	}{Items: searchData, Total: len(searchData)}}, nil
}
