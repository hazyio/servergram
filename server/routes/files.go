package routes

import (
	"hazyio.com/servergram/internal"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func RegisterFilesRoutes(server *internal.Server) {
	internal.Get(server, "/api/service/health", serverHealth)

}

func serverHealth(data internal.NormalRouteData[struct{}]) (*MessageResponse, error) {

	return &MessageResponse{Message: "all goo"}, nil
}
