package routes

import (
	"net/http"

	"gorm.io/gorm"
	"hazyio.com/servergram/server/internal"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func registerFilesRoutes(api *huma.Group, database *gorm.DB) {
	server := internal.NewGroup("/files", api, database, []string{"Files"})

	server.Get("/health", http.StatusOK, serverHealth)

}

func serverHealth(data internal.RouteData[struct{}]) (*EmptyRR, error) {

	return &EmptyRR{}, nil
}
