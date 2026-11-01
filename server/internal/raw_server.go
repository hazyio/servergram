package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal/sshmanager"
)

type RawServer struct {
	database       *gorm.DB
	server         *http.ServeMux
	allowedOrigins []string
	SSHManager     *sshmanager.SSHManager
	DataQueue      *DataQueue
	Polls          *Polls
}
type RawServerRouteData[I any] struct {
	Database   *gorm.DB
	Input      I
	DataQueue  *DataQueue
	Polls      *Polls
	SSHManager *sshmanager.SSHManager
}
type OnStartedRawServerRouteData struct {
	Database   *gorm.DB
	DataQueue  *DataQueue
	Polls      *Polls
	SSHManager *sshmanager.SSHManager
}

func NewRawServer(database *gorm.DB, mux *http.ServeMux, allowedOrigins []string, sshManager *sshmanager.SSHManager, polls *Polls) *RawServer {
	return &RawServer{
		database:       database,
		server:         mux,
		allowedOrigins: allowedOrigins,
		SSHManager:     sshManager,
		DataQueue:      NewDataQueue(),
		Polls:          polls,
	}
}

// WSHandler defines the type for WebSocket handlers, taking an input of type I wrapped in RawServerRouteData. close data queue if you want the connection to be closed gracefully.
func (s *RawServer) WebSocket[I any](pattern string, handle func(in RawServerRouteData[I]), runOnStart func(in OnStartedRawServerRouteData)) {
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
		dataQueue := s.DataQueue
		// watch for new messages in the data queue and send them over the WebSocket connection
		poller := s.Polls.AddPoll(NewPoller(time.Second, func(ctx context.Context) {
			messages, ok := dataQueue.Dequeue()
			if !ok {
				return
			}
			if messages == nil && dataQueue.Closed {
				// If there are no more messages and the queue is closed, close the connection.
				conn.Close(websocket.StatusNormalClosure, "no more messages, closing connection")
				return
			}
			for _, message := range messages {
				wsjson.Write(ctx, conn, message)
			}
		}))

		poller.Start(ctx)
		defer poller.Stop()
		// Run the function provided in runOnStart with the OnStartedRawServerRouteData.
		go runOnStart(OnStartedRawServerRouteData{
			Database:   s.database,
			DataQueue:  s.DataQueue,
			Polls:      s.Polls,
			SSHManager: s.SSHManager,
		})
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
			// Handle the incoming WebSocket message asynchronously.
			go handle(RawServerRouteData[I]{
				Database:   s.database,
				Input:      input,
				DataQueue:  s.DataQueue,
				Polls:      s.Polls,
				SSHManager: s.SSHManager,
			})
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
