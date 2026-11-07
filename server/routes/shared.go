package routes

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
)

func RegisterRoutes(api *huma.Group, database *gorm.DB, sshManager *sshmanager.SSHManager, allowedOrigins []string, polls *internal.Polls) {
	registerFilesRoutes(api, database, sshManager, allowedOrigins, polls)
	registerAccountRoutes(api, database, allowedOrigins, polls)
	registerSSHRoutes(api, database, allowedOrigins, polls)
	registerSearchRoutes(api, database, allowedOrigins, polls)
	registerGroupRoutes(api, database, allowedOrigins, polls)
	registerServerRoutes(api, database, sshManager, allowedOrigins, polls)
	registerLogRoutes(api, database, allowedOrigins, polls)
	registerNotificationRoutes(api, database, allowedOrigins, polls)
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
