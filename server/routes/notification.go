package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/routes/models"
)

func registerNotificationRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/notification", api, database, []string{"Notifications"}, nil)
	server.Get("/available", http.StatusOK, notificationAvailable)
	server.Get("/list/unread", http.StatusOK, unreadNotifications)
	server.Get("/list/read", http.StatusOK, readNotifications)
	server.Put("/mark-as-read", http.StatusOK, markAsRead)
	server.Get("/list/all", http.StatusOK, allNotifications)

}

func notificationAvailable(data internal.RouteData[EmptyRR]) (*BooleanResponse, error) {
	database := data.Database
	available, err := db.NewNotificationAvailable(database)
	if err != nil {
		return nil, err
	}
	response := &BooleanResponse{}
	response.Body.Value = available
	return response, nil
}

func unreadNotifications(data internal.RouteData[EmptyRR]) (*models.NotificationListResponse, error) {
	database := data.Database
	response := &models.NotificationListResponse{}
	response.Body = make([]models.NotificationListItem, 0)
	database.Model(&db.Notification{}).
		Select("notifications.id, notifications.title, notifications.message, notifications.url, notifications.created_at").
		Where("read = ?", false).
		Scan(&response.Body)
	return response, nil
}

func readNotifications(data internal.RouteData[EmptyRR]) (*models.NotificationListResponse, error) {
	database := data.Database
	response := &models.NotificationListResponse{}
	response.Body = make([]models.NotificationListItem, 0)
	database.Model(&db.Notification{}).
		Select("notifications.id, notifications.title, notifications.message, notifications.url, notifications.created_at").
		Where("read = ?", true).
		Scan(&response.Body)
	return response, nil
}

func allNotifications(data internal.RouteData[EmptyRR]) (*models.NotificationListResponse, error) {
	database := data.Database
	response := &models.NotificationListResponse{}
	response.Body = make([]models.NotificationListItem, 0)
	database.Model(&db.Notification{}).
		Select("notifications.id, notifications.title, notifications.message, notifications.url, notifications.created_at").
		Scan(&response.Body)
	return response, nil
}

func markAsRead(data internal.RouteData[models.MarkNotificationAsReadRequest]) (*BooleanResponse, error) {
	database := data.Database
	ids := data.Input.Body.NotificationIDs
	err := db.MarkNotificationsAsRead(database, ids)
	if err != nil {
		return nil, err
	}
	response := &BooleanResponse{}
	response.Body.Value = true
	return response, nil
}
