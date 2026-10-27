package routes

import (
	"gorm.io/gorm"
	"hazyio.com/servergram/internal"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func RegisterFilesRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/files", api, database)

	server.Get("/health", serverHealth)

}

func serverHealth(data internal.NormalRouteData[struct{}]) (*MessageResponse, error) {

	return &MessageResponse{Message: "all goo"}, nil
}
