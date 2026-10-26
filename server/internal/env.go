package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Env struct {
	DatabaseUrl    string   `validate:"required,url"`
	ServerPort     string   `validate:"required"`
	AllowedOrigins []string `validate:"required,dive,required"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {

		log.Fatal("Error loading .env file, are you sure it is in same dir you are running this application")
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
		log.Println("SERVER_PORT is not set, defaulting to 8080 this is dangerous as it allowed requests from all origins")
		env.ServerPort = "8080"
		env.AllowedOrigins = []string{"*"}
	}

	if err := validate.Struct(env); err != nil {
		log.Fatalf("environment validation failed: %v", err)
	}

	return env
}
