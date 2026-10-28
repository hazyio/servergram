package routes

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
)

func registerSearchRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/search", api, database)
	server.Post("/ssh", http.StatusOK, searchSSHKeys)
	server.Post("/group", http.StatusOK, searchGroups)
}

func searchSSHKeys(data internal.NormalRouteData[SearchRequest]) (*SearchResponse[SSHKeySearchData], error) {
	input := data.Input
	query := input.Body.Query
	database := data.Database

	var results []db.SSHKey

	if err := database.Where("name LIKE ?", "%"+query+"%").Limit(10).Find(&results).Error; err != nil {
		slog.Error("An error occurred while searching SSH keys", "error", err)
		return nil, internal.InternalServerError
	}

	var searchData []SSHKeySearchData
	for _, r := range results {
		searchData = append(searchData, SSHKeySearchData{
			Id:        r.ID,
			Name:      r.Name,
			Encrypted: r.Encrypted,
		})
	}

	return &SearchResponse[SSHKeySearchData]{Body: struct {
		Items []SSHKeySearchData `json:"items" required:"true" doc:"The search results"`
		Total int                `json:"total" required:"true" doc:"The total number of search results"`
	}{Items: searchData, Total: len(searchData)}}, nil
}

func searchGroups(data internal.NormalRouteData[SearchRequest]) (*SearchResponse[GroupSearchData], error) {
	input := data.Input
	query := input.Body.Query
	database := data.Database

	var results []db.Group

	if err := database.Where("name LIKE ?", "%"+query+"%").Limit(10).Find(&results).Error; err != nil {
		slog.Error("An error occurred while searching groups", "error", err)
		return nil, internal.InternalServerError
	}

	var searchData []GroupSearchData
	for _, r := range results {
		searchData = append(searchData, GroupSearchData{
			Id:   r.ID,
			Name: r.Name,
		})
	}

	return &SearchResponse[GroupSearchData]{Body: struct {
		Items []GroupSearchData `json:"items" required:"true" doc:"The search results"`
		Total int               `json:"total" required:"true" doc:"The total number of search results"`
	}{Items: searchData, Total: len(searchData)}}, nil
}
