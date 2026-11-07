package internal

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/sse"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal/sshmanager"

	"gorm.io/gorm"
)

type RouteData[I any] struct {
	Database   *gorm.DB
	Context    context.Context
	Input      I
	SSHManager *sshmanager.SSHManager
}

// Middleware wraps an http.Handler with additional behaviour.
type Middleware func(http.Handler) http.Handler

// Server registers routes on a ServeMux and exposes the database queries to them.
type Server struct {
	database       *gorm.DB
	group          *huma.Group
	mux            *http.ServeMux
	prefix         string
	sshManager     *sshmanager.SSHManager
	Polls          *Polls
	allowedOrigins []string
}

type SSEError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewSSEError(status int, message string) SSEError {
	return SSEError{Status: status, Message: message}
}
func NewSSEInternalError(message string) SSEError {
	return SSEError{Status: http.StatusInternalServerError, Message: message}
}

func SendSSEError(send sse.Sender, err error) {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.NotFound:
		_ = send.Data(NewSSEError(http.StatusNotFound, st.Message()))
	case codes.PermissionDenied:
		_ = send.Data(NewSSEError(http.StatusForbidden, st.Message()))
	case codes.OutOfRange:
		_ = send.Data(NewSSEError(http.StatusRequestEntityTooLarge, st.Message()))
	case codes.FailedPrecondition:
		_ = send.Data(NewSSEError(http.StatusBadRequest, st.Message()))
	default:
		_ = send.Data(NewSSEError(http.StatusInternalServerError, "Internal server error"))
	}
}

func FromGrpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		// not a gRPC status error
		return err
	}
	switch st.Code() {
	case codes.NotFound:
		return NewHumaError(http.StatusNotFound, st.Message())
	case codes.PermissionDenied:

		return NewHumaError(http.StatusForbidden, st.Message())
	case codes.AlreadyExists:
		return NewHumaError(http.StatusConflict, st.Message())
	case codes.InvalidArgument:
		return NewHumaError(http.StatusBadRequest, st.Message())
	default:
		return NewHumaError(http.StatusInternalServerError, "Internal server error")
	}
}

type bodyKey struct{}

// NewGroup wraps mux with a Server whose queries use pool.
func NewGroup(prefix string, api *huma.Group, database *gorm.DB, tags []string, sshManager *sshmanager.SSHManager, allowedOrigins []string, polls *Polls) *Server {
	group := huma.NewGroup(api, prefix)
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = tags
	})

	return &Server{
		database:       database,
		group:          group,
		prefix:         prefix,
		sshManager:     sshManager,
		Polls:          polls,
		allowedOrigins: allowedOrigins,
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
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})
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
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}

func (s *Server) Patch[I, O any](pattern string, status int, next func(RouteData[I]) (*O, error), errs ...int) {
	huma.Register(s.group, huma.Operation{
		OperationID:   "patch-" + opSlug(s.prefix, pattern),
		Summary:       huma.GenerateSummary(http.MethodPatch, pattern, new(O)),
		Method:        http.MethodPatch,
		Path:          pattern,
		DefaultStatus: status,
		Errors:        errs,
	}, func(ctx context.Context, input *I) (*O, error) {
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})
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
		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})
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

		resp, err := next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})
		if err != nil {
			return nil, ToHumaError(err)
		}
		return resp, nil
	})
}
func (s *Server) SSE[I, O any](pattern string, status int, next func(routeData RouteData[I], send sse.Sender) *O, errs ...int) {
	sse.Register(s.group, huma.Operation{
		OperationID: "sse-" + opSlug(s.prefix, pattern),
		Summary:     huma.GenerateSummary(http.MethodGet, pattern, new(O)),
		Method:      http.MethodGet,
		Path:        pattern,
		Errors:      errs,
	}, map[string]any{
		"data":  *new(O),
		"error": SSEError{},
	}, func(ctx context.Context, input *I, send sse.Sender) {
		next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager}, send)
	})
}
func (s *Server) Download[I any](pattern string, next func(routeData RouteData[I]) (*huma.StreamResponse, error), errs ...int) {

	huma.Register(s.group, huma.Operation{
		OperationID: "download-" + opSlug(s.prefix, pattern),
		Summary:     huma.GenerateSummary(http.MethodGet, pattern, struct{}{}),
		Method:      http.MethodGet,
		Path:        pattern,
		Errors:      errs,
		Responses: map[string]*huma.Response{
			"200": {
				Description: "File content",
				Content: map[string]*huma.MediaType{
					"application/octet-stream": {},
				},
			},
		},
	}, func(ctx context.Context, input *I) (*huma.StreamResponse, error) {
		return next(RouteData[I]{Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager})

	})
}
func (s *Server) Upload[I, O any](pattern string,
	maxBytes int64,
	next func(rd RouteData[I], body io.Reader) (*O, error),
	errs ...int) {
	type bodyKey struct{}

	// Runs before Huma parses input. Grabs the raw body without reading it.
	captureBody := func(hctx huma.Context, nextMw func(huma.Context)) {
		r, w := humago.Unwrap(hctx)

		// Huma's MaxBodyBytes only applies when Huma reads the body, so cap it here.
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

		// The server's ReadTimeout would kill long uploads; clear the deadline.
		_ = http.NewResponseController(w).SetReadDeadline(time.Time{})

		nextMw(huma.WithValue(hctx, bodyKey{}, io.Reader(r.Body)))
	}
	huma.Register(s.group, huma.Operation{
		OperationID:     "upload-" + opSlug(s.prefix, pattern),
		Summary:         huma.GenerateSummary(http.MethodPost, pattern, new(O)),
		Method:          http.MethodPost,
		Path:            pattern,
		MaxBodyBytes:    maxBytes,
		Errors:          errs,
		BodyReadTimeout: -1, // disables body read timeout
		Middlewares:     huma.Middlewares{captureBody},
		// Docs only: Huma won't parse this because I has no body field.
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				"application/octet-stream": {Schema: &huma.Schema{Type: "string", Format: "binary"}},
			},
		},
	}, func(ctx context.Context, input *I) (*O, error) {
		body, ok := ctx.Value(bodyKey{}).(io.Reader)
		if !ok {
			return nil, huma.Error500InternalServerError("body unavailable")
		}
		return next(RouteData[I]{
			Input: *input, Database: s.database, Context: ctx, SSHManager: s.sshManager,
		}, body)

	})
}

func (s *Server) WebSocket[I any](pattern string, next func(routeData RouteData[I], dataQueue *DataQueue), runOnStart func(in RouteData[any], dataQueue *DataQueue), errs ...int) {
	huma.Register(s.group, huma.Operation{
		OperationID:     "websocket-" + opSlug(s.prefix, pattern),
		Summary:         huma.GenerateSummary(http.MethodGet, pattern, struct{}{}),
		Method:          http.MethodGet,
		Path:            pattern,
		Errors:          errs,
		BodyReadTimeout: -1,
		Responses: map[string]*huma.Response{
			"101": {Description: "Switching Protocols"},
		},
	}, func(ctx context.Context, input *I) (*huma.StreamResponse, error) {
		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				r, w := humago.Unwrap(hctx)

				conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
					OriginPatterns: s.allowedOrigins,
				})
				if err != nil {
					return // Accept already wrote the error response
				}
				defer conn.CloseNow()
				conn.SetReadLimit(1 << 20)
				dataQueue := NewDataQueue()
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
				defer dataQueue.Close()
				defer poller.Stop()
				// Run the function provided in runOnStart with the OnStartedRawServerRouteData.
				go runOnStart(RouteData[any]{Input: &struct{}{}, Database: s.database, Context: ctx, SSHManager: s.sshManager}, dataQueue)
				for {
					var input I
					if err := wsjson.Read(ctx, conn, &input); err != nil {
						switch websocket.CloseStatus(err) {
						case websocket.StatusNormalClosure, websocket.StatusGoingAway:
						default:
							slog.Error("ws read error", "error", err)
						}
						poller.Stop()
						return
					}
					rd := RouteData[I]{Input: input, Database: s.database, Context: ctx, SSHManager: s.sshManager}
					select {
					case <-ctx.Done():
						return
					default:
						next(rd, dataQueue)

					}
				}
			},
		}, nil

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
