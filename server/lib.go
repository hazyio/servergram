package main

import (
	"context"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"

	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
	"hazyio.com/servergram/server/routes"
	"hazyio.com/servergram/server/websockets"
)

//go:embed all:dist
var distFS embed.FS

func NewRouter(_ context.Context, database *gorm.DB, polls *internal.Polls) *http.Server {
	env := internal.LoadEnv()
	cert, err := tls.LoadX509KeyPair(env.SSLPublicKey, env.SSLPrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr: ":" + env.ServerPort,
		// wrap the mux with CORS middleware
		Handler: internal.CorsMiddleware(env.AllowedOrigins, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
			MinVersion: tls.VersionTLS12,
		},
	}
	ssh_manager := sshmanager.NewSSHManager()
	config := huma.DefaultConfig("Sever Gram Api", "1.0.0")
	config.CreateHooks = nil // drops the "$schema" field, whose scheme/host is derived from the request
	api := humago.New(mux, config)
	apiUmbrella := huma.NewGroup(api, "/api")
	websocketGroup := huma.NewGroup(api, "")

	loggingMW := internal.LoggingMiddleware(api)
	sessionMW := internal.SessionMiddleware(database, api)
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		loggingMW(ctx, func(ctx huma.Context) {
			if skipAuthForRoutes(ctx.URL().Path) {
				slog.Info("Skipping authentication for route", "path", ctx.URL().Path)
				next(ctx)
				return
			}
			if env.Password == "" {
				slog.Info("Skipping authentication because password is not set", "path", ctx.URL().Path)
				next(ctx)
				return
			}
			// Middleware responses are invisible to Huma, so the 401 is added to the spec here.
			errSchema := api.OpenAPI().Components.Schemas.Schema(reflect.TypeFor[huma.ErrorModel](), true, "UnauthorizedErrorModel")
			apiUmbrella.UseSimpleModifier(func(op *huma.Operation) {
				if op.Responses == nil {
					op.Responses = map[string]*huma.Response{}
				}
				op.Responses["401"] = &huma.Response{
					Description: "Unauthorized",
					Content: map[string]*huma.MediaType{
						"application/problem+json": {Schema: errSchema},
					},
				}
			})
			sessionMW(ctx, next)

		})

	})
	// Register application routes and websockets
	routes.RegisterRoutes(apiUmbrella, database, ssh_manager, env.AllowedOrigins, polls)
	websockets.RegisterWebSockets(websocketGroup, database, ssh_manager, env.AllowedOrigins, polls)

	// handle root and frontend route
	mux.HandleFunc("/", handleRoot)

	// Write the OpenAPI specification to a file for reference
	b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
	_ = os.WriteFile("../openapi.json", b, 0644)
	return srv
}


func skipAuthForRoutes(path string) bool {
	// Add routes that should skip authentication here
	skipRoutes := []string{
		"/",                  // root route
		"/api/account/login", // login route
		"/_app/*",            // frontend route and its subroutes
	}
	for _, route := range skipRoutes {
		if strings.HasSuffix(route, "/*") {
			prefix := strings.TrimSuffix(route, "*")
			if path == strings.TrimSuffix(prefix, "/") || strings.HasPrefix(path, prefix) {
				return true
			}
		} else if path == route {
			return true
		}
	}
	return false
}
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"error":"endpoint not found"}`, http.StatusNotFound)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	sub, err := fs.Sub(distFS, "dist") // strip the "dist/" prefix
	if err != nil {
		log.Fatal(err)
	}
	http.FileServer(http.FS(sub)).ServeHTTP(w, r)

}

func InitDatabase() *gorm.DB {
	env := internal.LoadEnv()

	db_url := env.DatabaseUrl
	if db_url == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	dialector, err := getDatabaseDialector(db_url)
	if err != nil {
		log.Fatal(err)
	}
	database, err := gorm.Open(dialector, &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}

	if database.Dialector.Name() == "sqlite" {
		// SQLite serializes writes: one connection avoids "database is locked"
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(0) // never expire; an expired in-memory DB would vanish
		// WAL mode and a busy timeout so concurrent readers and a writer don’t collide
		database.Exec("PRAGMA journal_mode = WAL")
		database.Exec("PRAGMA busy_timeout = 5000")
		database.Exec("PRAGMA foreign_keys = ON")
	} else {
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	if err := database.AutoMigrate(db.Server{}, db.Sessions{}, db.SSHKey{}, db.Group{}, db.Log{}, db.Notification{}); err != nil {
		log.Fatal(err)
	}

	return database
}

func getDatabaseDialector(dbURL string) (gorm.Dialector, error) {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}
	switch strings.ToLower(parsedURL.Scheme) {
	case "sqlite", "sqlite3":
		dsn := strings.TrimPrefix(dbURL, parsedURL.Scheme+":")
		dsn = strings.TrimPrefix(dsn, "//")
		return sqlite.Open(dsn), nil
	case "postgres", "postgresql":
		return postgres.Open(dbURL), nil
	case "mysql":
		pw, _ := parsedURL.User.Password()
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
			parsedURL.User.Username(), pw, parsedURL.Host,
			strings.TrimPrefix(parsedURL.Path, "/"), parsedURL.RawQuery)
		if parsedURL.RawQuery == "" {
			dsn = strings.TrimSuffix(dsn, "?") + "?parseTime=true"
		}
		return mysql.Open(dsn), nil
	case "sqlserver", "mssql":
		return sqlserver.Open(dbURL), nil
	default:
		return nil, fmt.Errorf("unsupported database scheme: %q", parsedURL.Scheme)
	}
}
