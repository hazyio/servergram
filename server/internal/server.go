package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"hazyio.com/servergram/db"

	"gorm.io/gorm"
)

type NormalRouteData struct {
	Request        *http.Request
	Response       *JsonResponse
	ResponseWriter http.ResponseWriter
	Database       *gorm.DB
}
type AuthenticatedRouteData struct {
	Request        *http.Request
	Response       *JsonResponse
	ResponseWriter http.ResponseWriter
	Database       *gorm.DB
	Service        db.Service
}

// NormalRouteHandler handles a request that needs no authentication.
type NormalRouteHandler func(NormalRouteData)

// AuthenticatedRouteHandler handles a request on behalf of an authenticated user.
type AuthenticatedRouteHandler func(AuthenticatedRouteData)

// Middleware wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Server registers routes on a ServeMux and exposes the database queries to them.
type Server struct {
	database *gorm.DB
	server   *http.ServeMux
}

// NewServer wraps mux with a Server whose queries use pool.
func NewServer(server *http.ServeMux, database *gorm.DB) *Server {
	return &Server{
		database: database,
		server:   server,
	}

}

// Route registers an unauthenticated handler wrapped with recovery and logging.
func (s *Server) Route(pattern string, next NormalRouteHandler) {
	s.server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(NormalRouteData{Request: r, Response: NewResponse(w), ResponseWriter: w, Database: s.database})
		})
		chain(h, recoverer, logging).ServeHTTP(w, r)
	})
}

// AuthenticatedRoute registers a handler that requires a well-formed (32-char) "muster_session" cookie,
// responding 401 otherwise. It does not yet look the session up in the database.
// It returns the Server's queries.
func (s *Server) AuthenticatedRoute(pattern string, next AuthenticatedRouteHandler) {
	// ctx := context.Background()
	s.server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		session_token, err := r.Cookie("session_token")
		response := NewResponse(w)
		if err != nil {
			response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, session token not found"})
			return
		}
		trimmed_session_token := strings.TrimSpace(session_token.Value)
		if trimmed_session_token == "" {
			response.Write(http.StatusUnauthorized, map[string]string{"error": "session token is empty"})
			return
		}

		// session, err := s.database.GetSessionByToken(ctx, trimmed_session_token)
		// var session db.Sessions

		// err = s.database.WithContext(ctx).First(&session, "token = ?", HashToken(trimmed_session_token)).Error
		// if err != nil {
		// 	response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, session not found"})
		// 	return
		// }
		// var user db.User
		// err = s.database.WithContext(ctx).First(&user, "id = ?", session.UserID).Error
		// if err != nil {
		// 	response.Write(http.StatusUnauthorized, map[string]string{"error": "unauthorized, user not found"})
		// 	return
		// }
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next(AuthenticatedRouteData{Request: r, Response: NewResponse(w), ResponseWriter: w, Service: db.Service{}, Database: s.database})
		})
		chain(h, recoverer, logging).ServeHTTP(w, r)
	})

}

// chain applies m to h, then recoverer and logging, so the first middleware is outermost.
func chain(h http.Handler, m ...Middleware) http.Handler {
	for _, v := range slices.Backward(m) {
		h = v(h)
	}
	return h
}

// recoverer turns a panic in next into a 500 JSON response.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err)
				response := NewResponse(w)
				response.Write(http.StatusInternalServerError, map[string]string{"error": "An internal error occurred"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// logging logs the method, path and duration of each request after it completes.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // call the wrapped handler
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

// func handleNotFound(w http.ResponseWriter, r *http.Request) {
// 	response := NewResponse(w)
// 	response.Write(http.StatusNotFound, map[string]string{"error": "endpoint not found"})
// }

// func handleRoot(w http.ResponseWriter, r *http.Request) {
// 	handleNotFound(w, r)
// }

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]) // 64 chars
}
