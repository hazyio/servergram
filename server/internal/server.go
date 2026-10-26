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

	// s.server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
	// 	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		next(NormalRouteData{Request: r, Response: NewResponse(w), ResponseWriter: w, Database: s.database})
	// 	})
	// 	chain(h, recoverer, logging).ServeHTTP(w, r)
	// })
}

// AuthenticatedRoute registers a handler that requires a well-formed (32-char) "muster_session" cookie,
// responding 401 otherwise. It does not yet look the session up in the database.
// It returns the Server's queries.
// func (s *Server) AuthenticatedRoute(pattern string, next AuthenticatedRouteHandler) {
// 	// ctx := context.Background()
// 	s.server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
// 		session_token, err := r.Cookie("session_token")
// 		response := NewResponse(w)
// 		if err != nil {
// 			response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, session token not found"})
// 			return
// 		}
// 		trimmed_session_token := strings.TrimSpace(session_token.Value)
// 		if trimmed_session_token == "" {
// 			response.Write(http.StatusUnauthorized, map[string]string{"error": "session token is empty"})
// 			return
// 		}

// 		// session, err := s.database.GetSessionByToken(ctx, trimmed_session_token)
// 		// var session db.Sessions

// 		// err = s.database.WithContext(ctx).First(&session, "token = ?", HashToken(trimmed_session_token)).Error
// 		// if err != nil {
// 		// 	response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, session not found"})
// 		// 	return
// 		// }
// 		// var user db.User
// 		// err = s.database.WithContext(ctx).First(&user, "id = ?", session.UserID).Error
// 		// if err != nil {
// 		// 	response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, user not found"})
// 		// 	return
// 		// }
// 		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			next(AuthenticatedRouteData{Request: r, Response: NewResponse(w), ResponseWriter: w, Service: db.Service{}, Database: s.database})
// 		})
// 		chain(h, recoverer, logging).ServeHTTP(w, r)
// 	})

// }
