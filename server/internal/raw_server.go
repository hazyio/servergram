package internal

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"
	"hazyio.com/servergram/db"
)

type RawServer struct {
	database *gorm.DB
	server   *http.ServeMux
	prefix   string
}
type RawServerRouteData struct {
	ResponseWriter http.ResponseWriter
	Database       *gorm.DB
	Service        db.Server
}

func (s *RawServer) WebSocket[I, O any](pattern string) {
	// h :=
	s.server.Handle("GET /ws"+pattern, WebSocketSessionMiddleware(s.database, func(w http.ResponseWriter, r *http.Request) {

		// conn, err := websocket.Accept(w, r, nil)
		// if err != nil {
		// 	return
		// }
		// defer conn.CloseNow()
		// read/write loop
		w.Write([]byte("hellow"))
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
