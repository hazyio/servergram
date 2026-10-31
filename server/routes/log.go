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

func listLogs(data internal.NormalRouteData[models.ListLogRequest]) (*models.LogListResponse, error) {
	database := data.Database
	input := data.Input
	response := &models.LogListResponse{}
	response.Body.Items = make([]models.LogListItem, 0)

	database.Model(&db.Log{}).
		Select("logs.id, logs.action, logs.message, logs.created_at").
		Where("logs.id > ?", input.StartIndex).Order("logs.id ASC").Limit(int(input.Limit)).
		Scan(&response.Body.Items)
	if len(response.Body.Items) > 0 {
		firstID := response.Body.Items[0].Id
		lastID := response.Body.Items[len(response.Body.Items)-1].Id

		// Set the next and previous index for pagination
		response.Body.NextIndex = lastID
		if int32(firstID)-int32(input.Limit) <= 0 {
			response.Body.PreviousIndex = 0
		} else {
			response.Body.PreviousIndex = firstID - input.Limit - 1
		}

	}

	return response, nil
}
