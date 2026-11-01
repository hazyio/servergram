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

func registerLogRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/logs", api, database, []string{"Logs"})
	server.Get("/list", http.StatusOK, listLogs)

}

func listLogs(data internal.NormalRouteData[models.ListLogRequest]) (*models.LogListResponse, error) {
	database := data.Database
	input := data.Input
	response := &models.LogListResponse{}
	response.Body.Items = make([]models.LogListItem, 0)

	query := database.Model(&db.Log{}).
		Select("logs.id, logs.action, logs.message, logs.created_at")

	if input.StartIndex != 0 {
		query = query.Where("logs.id < ?", input.StartIndex)
	}
	query.Order("logs.created_at DESC").Limit(int(input.Limit)).
		Scan(&response.Body.Items)
	if len(response.Body.Items) > 0 {
		firstID := response.Body.Items[0].Id
		lastID := response.Body.Items[len(response.Body.Items)-1].Id

		response.Body.NextIndex = lastID
		if input.StartIndex != 0 {
			newPreviousIndex := firstID + input.Limit + 1
			var remaining int64
			err := database.Model(&db.Log{}).
				Select("id").
				Where("id < ?", newPreviousIndex).
				Where("id > ?", newPreviousIndex-input.Limit). // Check for logs within the previous range
				Count(&remaining).Error
			slog.Info("Remaining previous logs", "remaining", remaining)
			if err != nil {
				slog.Error("Error checking for previous logs", "err", err)
			} else if remaining == 0 {
				// if there are fewer previous logs than the limit, set PreviousIndex to 0
				response.Body.PreviousIndex = 0

			} else {
				// if there is an index available, use it as the previous index
				response.Body.PreviousIndex = newPreviousIndex
			}
		}

	}

	return response, nil
}
