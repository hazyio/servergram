package routes

import (
	"net/http"

	"hazyio.com/servergram/internal"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func HandleServer(server *internal.Server) {
	server.Route("GET /api/service/health", serverHealth)

}

func serverHealth(data internal.NormalRouteData) {
	data.Response.Write(http.StatusOK, map[string]string{"message": "all good"})
}
