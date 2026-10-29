package routes

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(api *huma.Group, database *gorm.DB) {
	registerFilesRoutes(api, database)
	registerAccountRoutes(api, database)
	registerSSHRoutes(api, database)
	registerSearchRoutes(api, database)
	registerGroupRoutes(api, database)
}

type EmptyRR struct {
}

type PathIdRequest struct {
	ID uint `path:"id" minimum:"1" doc:"The ID of the current resource"`
}
