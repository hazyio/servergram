package routes

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/internal/sshmanager"
)

func RegisterRoutes(api *huma.Group, database *gorm.DB, sshManager *sshmanager.SSHManager) {
	registerFilesRoutes(api, database, sshManager)
	registerAccountRoutes(api, database)
	registerSSHRoutes(api, database)
	registerSearchRoutes(api, database)
	registerGroupRoutes(api, database)
	registerServerRoutes(api, database, sshManager)
	registerLogRoutes(api, database)
	registerNotificationRoutes(api, database)
}

type EmptyRR struct {
}
type BodyResponse struct {
	Body any
}

type PathIdRequest struct {
	ID uint `path:"id" minimum:"1" doc:"The ID of the current resource"`
}

type BooleanResponse struct {
	Body struct {
		Value bool `json:"value"`
	}
}
