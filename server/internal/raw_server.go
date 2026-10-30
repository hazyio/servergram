package internal

import (
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"gorm.io/gorm"
	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal/sshmanager"
)

type RawServer struct {
	database       *gorm.DB
	server         *http.ServeMux
	allowedOrigins []string
	SSHManager     *sshmanager.SSHManager
}
type RawServerRouteData[I any] struct {
	Database *gorm.DB
	Input    I
}

func NewRawServer(database *gorm.DB, mux *http.ServeMux, allowedOrigins []string, sshManager *sshmanager.SSHManager) *RawServer {
	return &RawServer{
		database:       database,
		server:         mux,
		allowedOrigins: allowedOrigins,
		SSHManager:     sshManager,
	}
}

// WSHandler defines the type for WebSocket handlers, taking an input of type I and returning an output of type O and an error. when you return an error, the socket connection will be closed.
func (s *RawServer) WebSocket[I any](pattern string, handle func(in RawServerRouteData[I]) (any, error)) {
	s.server.Handle("GET /ws"+pattern, WebSocketSessionMiddleware(s.database, func(w http.ResponseWriter, req *http.Request) {

		conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
			OriginPatterns: s.allowedOrigins,
		})
		if err != nil {
			slog.Error("ws accept error", "error", err)

			return // Accept already responded
		}
		defer conn.CloseNow()
		ctx := req.Context()

		for {
			c, err := req.Cookie("session_token")
			if err != nil {
				conn.Close(websocket.StatusPolicyViolation, "missing session token")
				return
			}
			if !db.TokenIsValid(HashString(c.Value), s.database) {
				conn.Close(websocket.StatusPolicyViolation, "invalid session token")
				return
			}
			var input I
			if err := wsjson.Read(ctx, conn, &input); err != nil {
				switch websocket.CloseStatus(err) {
				case websocket.StatusNormalClosure, websocket.StatusGoingAway:
				default:
					slog.Error("ws read error", "error", err)
				}
				return
			}
			output, err := handle(RawServerRouteData[I]{
				Database: s.database,
				Input:    input,
			})
			if err != nil {
				conn.Close(websocket.StatusInternalError, "function errored, closing connection")
				return
			}
			if output != nil {
				wsjson.Write(ctx, conn, output)
			}
		}

	}))
}

func WebSocketSessionMiddleware(database *gorm.DB, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Implement session checking logic here
		c, err := req.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if c.Value == "" {
			slog.Warn("Empty session token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !db.TokenIsValid(HashString(c.Value), database) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	})
}
