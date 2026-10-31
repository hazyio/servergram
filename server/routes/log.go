package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/routes/models"
)

func registerLogRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/logs", api, database, []string{"Logs"})
	server.Get("/list", http.StatusOK, listLogs)

}

func listLogs(data internal.NormalRouteData[*models.ListLogRequest]) (*models.LogListResponse, error) {
	database := data.Database
	input := data.Input.Body
	response := &models.LogListResponse{}
	response.Body.Items = make([]models.LogListItem, 0)

	database.Model(&db.Log{}).
		Select("logs.id, logs.action, logs.message, logs.created_at").
		Where("logs.id > ?", input.StartIndex).
		Order("logs.id ASC").
		Limit(int(input.Collect)).
		Scan(&response.Body.Items)

	return response, nil
}
