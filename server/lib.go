package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"

	"hazyio.com/servergram/db"
	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/routes"
)

//go:embed all:dist
var distFS embed.FS

func NewRouter(database *gorm.DB) http.Handler {
	env := internal.LoadEnv()
	loginTokens := make(internal.LoginToken)

	mux := http.NewServeMux()
	slog.Info("Logger initialized")
	config := huma.DefaultConfig("Sever Gram Api", "1.0.0")
	config.CreateHooks = nil // drops the "$schema" field, whose scheme/host is derived from the request

	api := humago.New(mux, config)

	server := internal.NewServer(api, database)
	if env.Password != "" {
		// add middleware only if password is set
		api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
			if skipAuthForRoutes(ctx.URL().Path) {
				next(ctx)
				return
			}
			internal.SessionMiddleware(&loginTokens, api)(ctx, next)
		})
	}

	routes.RegisterFilesRoutes(server)
	mux.HandleFunc("/", handleRoot)

	return internal.CORS(env.AllowedOrigins)(mux)
}
func skipAuthForRoutes(path string) bool {
	// Add routes that should skip authentication here
	skipRoutes := []string{
		"/",        // root route
		"/login/*", // login route and its subroutes
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
	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := database.AutoMigrate(db.Service{}); err != nil {
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
		return mysql.Open(dbURL), nil
	case "sqlserver", "mssql":
		return sqlserver.Open(dbURL), nil
	default:
		return nil, fmt.Errorf("unsupported database scheme: %q", parsedURL.Scheme)
	}
}
