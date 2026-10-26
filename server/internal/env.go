package internal

import (
	"net/url"
	"strings"

	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseUrl    string   `validate:"required,database_url"`
	ServerPort     string   `validate:"required"`
	AllowedOrigins []string `validate:"required,dive,required"`
}

var validate = func() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("database_url", func(field validator.FieldLevel) bool {
		parsedURL, err := url.Parse(field.Field().String())
		if err != nil {
			return false
		}

		switch strings.ToLower(parsedURL.Scheme) {
		case "sqlite", "sqlite3":
			return parsedURL.Opaque != "" || parsedURL.Host != "" || parsedURL.Path != ""
		case "postgres", "postgresql", "mysql", "sqlserver", "mssql":
			return parsedURL.Host != ""
		default:
			return false
		}
	})
	return v
}()

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Error loading .env file, are you sure it is in same dir you are running this application")
		}
	}

	env := Env{

		DatabaseUrl:    os.Getenv("DATABASE_URL"),
		ServerPort:     os.Getenv("SERVER_PORT"),
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
	}
	for i, origin := range env.AllowedOrigins {
		env.AllowedOrigins[i] = strings.TrimSpace(origin)
	}

	if strings.TrimSpace(env.DatabaseUrl) == "" {
		// If DATABASE_URL is not set, default to a SQLite database in the config directory.
		dir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("failed to get user config dir: %v", err)
		}
		env.DatabaseUrl = "sqlite3://" + dir + "/servergram.db"
	}

	if strings.TrimSpace(env.ServerPort) == "" {
		log.Println("SERVER_PORT is not set, defaulting to 9788 this is dangerous as it allowed requests from all origins")
		env.ServerPort = "9788"
		env.AllowedOrigins = []string{"*"}
	}
	println("Loaded environment variables: ", env.DatabaseUrl)

	if err := validate.Struct(env); err != nil {
		log.Fatalf("environment validation failed: %v", err)
	}

	return env
}
