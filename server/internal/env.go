package internal

import (
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseUrl    string   `validate:"required,database_url"`
	ServerPort     string   `validate:"required"`
	AllowedOrigins []string `validate:"required,dive,required"`
	Password       string   `validate:"omitempty,strong_password"`
	ShhMasterKey   string   `validate:"required,strong_password"`
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
	_ = v.RegisterValidation("strong_password", func(field validator.FieldLevel) bool {
		password := field.Field().String()
		length := utf8.RuneCountInString(password)
		if length < 16 || length > 128 {
			return false
		}

		var hasUpper, hasLower, hasNumber, hasSymbol bool
		for _, char := range password {
			if !unicode.IsPrint(char) || unicode.IsSpace(char) {
				return false
			}
			hasUpper = hasUpper || unicode.IsUpper(char)
			hasLower = hasLower || unicode.IsLower(char)
			hasNumber = hasNumber || unicode.IsNumber(char)
			hasSymbol = hasSymbol || unicode.IsPunct(char) || unicode.IsSymbol(char)
		}
		return hasUpper && hasLower && hasNumber && hasSymbol
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
		Password:       os.Getenv("PASSWORD"),
		ShhMasterKey:   os.Getenv("SHH_MASTER_KEY"),
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
	if env.Password == "" {
		log.Println("PASSWORD is not set, authentication will be disabled")
	}

	if strings.TrimSpace(env.ServerPort) == "" {
		log.Println("SERVER_PORT is not set, defaulting to 9788 this is dangerous as it allowed requests from all origins")
		env.ServerPort = "9788"
		env.AllowedOrigins = []string{"*"} 
	}
	if env.Password == env.ShhMasterKey {
		log.Fatalf("PASSWORD and SHH_MASTER_KEY should not be the same")
	}
	if err := validate.Struct(env); err != nil {
		log.Fatalf("environment validation failed: %v", err)
	}

	return env
}
