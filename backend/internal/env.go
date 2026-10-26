package internal

import (
	"log"

	"github.com/joho/godotenv"
)

type Env struct {
	ResendAPIKey string
	DatabaseUrl  string
	ServerPort   string
	AppUrl       string
}

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {

		log.Fatal("Error loading .env file, are you sure it is in same dir you are running this application")
	}

	env := Env{}

	// missing := make([]string, 0, 4)
	// if strings.TrimSpace(env.ResendAPIKey) == "" {
	// 	missing = append(missing, "RESEND_API_KEY")
	// }
	// if strings.TrimSpace(env.DatabaseUrl) == "" {
	// 	missing = append(missing, "DATABASE_URL")
	// }
	// if strings.TrimSpace(env.ServerPort) == "" {
	// 	missing = append(missing, "SERVER_PORT")
	// }
	// if strings.TrimSpace(env.AppUrl) == "" {
	// 	missing = append(missing, "APP_URL")
	// }
	// if len(missing) > 0 {
	// 	log.Fatalf("required environment variable(s) are missing or blank: %s", strings.Join(missing, ", "))
	// }

	return env
}
