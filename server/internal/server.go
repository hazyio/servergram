package internal

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal/sshmanager"

	"gorm.io/gorm"
)

type RouteData[I any] struct {
	Database *gorm.DB
	Context  *context.Context
	Input    I
}

// Middleware wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Server registers routes on a ServeMux and exposes the database queries to them.
type Server struct {
	database   *gorm.DB
	group      *huma.Group
	mux        *http.ServeMux
	prefix     string
	sshManager *sshmanager.SSHManager
}

// NewGroup wraps mux with a Server whose queries use pool.
func NewGroup(prefix string, api *huma.Group, database *gorm.DB, tags []string, sshManager *sshmanager.SSHManager) *Server {
	group := huma.NewGroup(api, prefix)
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = tags
	})

	return &Server{
		database:   database,
		group:      group,
		prefix:     prefix,
		sshManager: sshManager,
	}

}

func opSlug(parts ...string) string {
	s := strings.ToLower(strings.Join(parts, "-"))
	s = strings.NewReplacer("/", "-", "{", "", "}", "").Replace(s)
	return strings.Trim(s, "-")
}
func (s *Server) Get[I, O any](pattern string, status int, next func(RouteData[I]) (*O, error), errs ...int) {

	huma.Register(s.group, huma.Operation{
		OperationID:   "get-" + opSlug(s.prefix, pattern),
		Summary:       huma.GenerateSummary(http.MethodGet, pattern, new(O)),
		Method:        http.MethodGet,
		Path:          pattern,
		DefaultStatus: status,
		Errors:        errs,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}

func (s *Server) Delete[I, O any](pattern string, status int, next func(RouteData[I]) (*O, error), errs ...int) {
	huma.Register(s.group, huma.Operation{
		OperationID:   "delete-" + opSlug(s.prefix, pattern),
		Summary:       huma.GenerateSummary(http.MethodDelete, pattern, new(O)),
		Method:        http.MethodDelete,
		Path:          pattern,
		DefaultStatus: status,
		Errors:        errs,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}
func (s *Server) Post[I, O any](pattern string, status int, next func(RouteData[I]) (*O, error), errs ...int) {
	huma.Register(s.group, huma.Operation{
		OperationID:   "post-" + opSlug(s.prefix, pattern),
		Summary:       huma.GenerateSummary(http.MethodPost, pattern, new(O)),
		Method:        http.MethodPost,
		Path:          pattern,
		DefaultStatus: status,
		Errors:        errs,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: &ctx})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}
func (s *Server) Put[I, O any](pattern string, status int, next func(RouteData[I]) (*O, error), errs ...int) {
	huma.Register(s.group, huma.Operation{
		OperationID:   "put-" + opSlug(s.prefix, pattern),
		Summary:       huma.GenerateSummary(http.MethodPut, pattern, new(O)),
		Method:        http.MethodPut,
		Path:          pattern,
		DefaultStatus: status,
		Errors:        errs,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: &ctx})
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
		if !db.TokenIsValid(HashString(c.Value), database) {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next(ctx)
	}

}
func CorsMiddleware(allowedOrigins []string, next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowed[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept-Language")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
