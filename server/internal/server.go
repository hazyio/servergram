package internal

import (
	"context"
	"net/http"

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

// AuthenticatedRouteHandler handles a request on behalf of an authenticated user.
type AuthenticatedRouteHandler func(AuthenticatedRouteData)

// Middleware wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Server registers routes on a ServeMux and exposes the database queries to them.
type Server struct {
	database *gorm.DB
	group    *huma.Group
}

// NewGroup wraps mux with a Server whose queries use pool.
func NewGroup(prefix string, api *huma.Group, database *gorm.DB) *Server {
	group := huma.NewGroup(api, prefix)

	return &Server{
		database: database,
		group:    group,
	}

}

// Get registers an unauthenticated GET handler wrapped with recovery and logging.
// Methods cannot have type parameters, so this is a package-level function.
func (s *Server) Get[I, O any](pattern string, next func(NormalRouteData[I]) (*O, error)) {
	huma.Get(s.group, pattern, func(
		ctx context.Context,
		input *I,
	) (*JsonResponse[O], error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		if resp == nil {
			return nil, huma.Error500InternalServerError("internal error")
		}

		return NewResponse(*resp), nil
	})

}

// Post registers an unauthenticated POST handler wrapped with recovery and logging.
// Methods cannot have type parameters, so this is a package-level function.
func (s *Server) Post[I, O any](pattern string, next func(NormalRouteData[I]) (*O, error)) {
	huma.Post(s.group, pattern, func(
		ctx context.Context,
		input *I,
	) (*JsonResponse[O], error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		if resp == nil {
			return nil, huma.Error500InternalServerError("internal error")
		}

		return NewResponse(*resp), nil
	})

}
func (s *Server) PostRaw[I, O any](pattern string, next func(NormalRouteData[I]) (*O, error)) {
	huma.Post(s.group, pattern, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
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
		if !loginTokens.Valid(c.Value) {
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
