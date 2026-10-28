package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"hazyio.com/servergram/db"

	"gorm.io/gorm"
)

type ctxKey struct{}

type NormalRouteData[I any] struct {
	Database *gorm.DB
	Context  *context.Context
	Input    I
}
type AuthenticatedRouteData struct {
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

func (s *Server) Get[I, O any](pattern string, status int, next func(NormalRouteData[I]) (*O, error)) {
	huma.Register(s.group, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodGet, pattern, new(O)),
		Summary:       huma.GenerateSummary(http.MethodGet, pattern, new(O)),
		Method:        http.MethodGet,
		Path:          pattern,
		DefaultStatus: status,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}

func (s *Server) Delete[I, O any](pattern string, status int, next func(NormalRouteData[I]) (*O, error)) {
	huma.Register(s.group, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodDelete, pattern, new(O)),
		Summary:       huma.GenerateSummary(http.MethodDelete, pattern, new(O)),
		Method:        http.MethodDelete,
		Path:          pattern,
		DefaultStatus: status,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}
func (s *Server) Post[I, O any](pattern string, status int, next func(NormalRouteData[I]) (*O, error)) {
	huma.Register(s.group, huma.Operation{
		OperationID:   huma.GenerateOperationID(http.MethodPost, pattern, new(O)),
		Summary:       huma.GenerateSummary(http.MethodPost, pattern, new(O)),
		Method:        http.MethodPost,
		Path:          pattern,
		DefaultStatus: status,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(NormalRouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}
func SessionMiddleware(database *gorm.DB, api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		req := http.Request{Header: http.Header{"Cookie": {ctx.Header("Cookie")}}}
		c, err := req.Cookie("session_token")
		if err != nil {

			if isNavigating(ctx) {

				return
			}
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if c.Value == "" {
			slog.Warn("Empty session token")
			if isNavigating(ctx) {

				return
			}
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		find_session := db.Sessions{}
		if err := database.Where("token = ?", HashString(c.Value)).First(&find_session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Warn("Session not found")
				if isNavigating(ctx) {

					return
				}
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
				return
			}
			slog.Error("An internal error occurred while trying to check the session", "error", err)
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if find_session.IsExpired() {
			slog.Warn("Session expired")
			if isNavigating(ctx) {
				return
			}
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next(ctx)
	}

}
func isNavigating(ctx huma.Context) bool {
	if ctx.Header("Sec-Fetch-Mode") == "navigate" {
		ctx.SetHeader("Location", "/account/login")
		ctx.SetStatus(http.StatusSeeOther)
		return true
	}
	return false
}
func LoggingMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		slog.Info("Request received", "method", ctx.Method(), "path", ctx.URL().Path)
		next(ctx)
	}
}
