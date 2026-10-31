package models

type MarkNotificationAsReadRequest struct {
	Body struct {
		NotificationIDs []int `json:"notification_ids" required:"true" doc:"The IDs of the notifications to mark as read"`
	}
}

type NotificationListItem struct {
	ID        int     `json:"id" doc:"The ID of the notification"`
	Title     string  `json:"title" doc:"The title of the notification"`
	Message   string  `json:"message" doc:"The message of the notification"`
	Url       *string `json:"url" doc:"The URL associated with the notification"`
	CreatedAt string  `json:"created_at" doc:"The creation timestamp of the notification"`
}

type NotificationListResponse struct {
	Body []NotificationListItem `json:"body" doc:"The list of notifications"`
}
