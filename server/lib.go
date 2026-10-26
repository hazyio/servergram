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

	mux := http.NewServeMux()
	slog.Info("Logger initialized")
	api := humago.New(mux, huma.DefaultConfig("Sever Gram Api", "1.0.0"))
	server := internal.NewServer(api, database)
	routes.RegisterFilesRoutes(server)
	mux.HandleFunc("/", handleRoot)

	return internal.CORS(env.AllowedOrigins)(mux)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	response := internal.NewResponse(w)
	response.Write(http.StatusNotFound, map[string]string{"error": "endpoint not found"})
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
