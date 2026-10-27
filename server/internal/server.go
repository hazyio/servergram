package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"hazyio.com/servergram/db"

	"gorm.io/gorm"
)

type NormalRouteData[I any] struct {
	Database *gorm.DB
	Context  *context.Context
	Input    I
}
type AuthenticatedRouteData struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Database       *gorm.DB
	Service        db.Service
}
type LoginTokenInformation struct {
	Expires time.Time
}
type LoginToken map[string]LoginTokenInformation

// AuthenticatedRouteHandler handles a request on behalf of an authenticated user.
type AuthenticatedRouteHandler func(AuthenticatedRouteData)

// Middleware wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Server registers routes on a ServeMux and exposes the database queries to them.
type Server struct {
	database *gorm.DB
	server   huma.API
}

// NewServer wraps mux with a Server whose queries use pool.
func NewServer(server huma.API, database *gorm.DB) *Server {
	return &Server{
		database: database,
		server:   server,
	}

}

// Get registers an unauthenticated GET handler wrapped with recovery and logging.
// Methods cannot have type parameters, so this is a package-level function.
func Get[I, O any](s *Server, pattern string, next func(NormalRouteData[I]) (*O, error)) {
	huma.Get(s.server, pattern, func(
		ctx context.Context,
		input *I,
	) (*JsonResponse[O], error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, err
		}
		return NewResponse(*resp), nil
	})

}
func SessionMiddleware(loginTokens *LoginToken, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		req := http.Request{Header: http.Header{"Cookie": {ctx.Header("Cookie")}}}
		c, err := req.Cookie("session_token")
		if err != nil || c.Value == "" {
			redirect(ctx, "/account/login")

			return
		}
		info, ok := (*loginTokens)[c.Value]
		if !ok || info.Expires.Before(time.Now()) {
			redirect(ctx, "/account/login")
			return
		}

		next(ctx)
	}

}

func redirect(ctx huma.Context, url string) {
	ctx.SetHeader("Location", url)
	ctx.SetStatus(http.StatusSeeOther) // 303
}
