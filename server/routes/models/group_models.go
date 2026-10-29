package models
import "time"
type CreateGroupRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" maxLength:"100" pattern:"^[A-Za-z0-9_ .@-]+$" doc:"The name of the group"`
	}
}
type GroupListResponse struct {
	Body struct {
		Items []struct {
			Id                uint      `json:"id"`
			Name              string    `json:"name"`
			CreatedAt         time.Time `json:"created_at"`
			ResourcesAttached uint      `json:"resources_attached"`
		} `json:"items" required:"true" doc:"The list of groups"`
		Total int `json:"total" required:"true" doc:"The total number of groups"`
	}
}
