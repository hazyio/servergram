package routes

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mileusna/useragent"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/routes/models"
)

func registerAccountRoutes(api *huma.Group, database *gorm.DB, allowedOrigins []string, polls *internal.Polls) {
	server := internal.NewGroup("/account", api, database, []string{"Accounts"}, nil, allowedOrigins, polls)
	server.Post("/login", http.StatusCreated, login)
	server.Delete("/logout", http.StatusOK, logOut)

}

func login(data internal.RouteData[models.LoginRequest]) (*models.LoginSuccessfulRedirectResponse, error) {
	input := data.Input
	password := input.Body.Password
	env := internal.LoadEnv()
	if env.Password == "" {
		slog.Error("server password is not configured")
		return nil, internal.NewHumaError(http.StatusInternalServerError, "server password is not configured")
	}
	a := sha256.Sum256([]byte(password))
	b := sha256.Sum256([]byte(env.Password))
	if subtle.ConstantTimeCompare(a[:], b[:]) != 1 {
		return nil, internal.NewHumaError(http.StatusUnauthorized, "Password incorrect")
	}
	randomToken, err := internal.GenerateRandomString(32)
	if err != nil {
		slog.Error("error generating random token", "err", err)
		return nil, internal.NewHumaError(http.StatusInternalServerError, "error generating random token")
	}
	expires := time.Now().Add(24 * time.Hour)
	ua := useragent.Parse(input.UserAgent)
	new_session := db.Sessions{
		Token:     internal.HashString(randomToken),
		Expires:   expires,
		Os:        ua.OS,
		UserAgent: input.UserAgent,
	}
	if err := data.Database.Create(&new_session).Error; err != nil {
		slog.Error("error creating new session", "err", err)
		return nil, internal.NewHumaError(http.StatusInternalServerError, "error creating new session")
	}
	if err := data.Database.Where("expires < ?", time.Now()).Delete(&db.Sessions{}).Error; err != nil {
		slog.Error("error deleting expired sessions", "err", err)
	}
	db.AddLog(db.LogActionCreate, "User logged in", data.Database)
	return &models.LoginSuccessfulRedirectResponse{
		Status: http.StatusCreated,
		SetCookie: http.Cookie{
			Name:     "session_token",
			Value:    randomToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Expires:  expires,
			MaxAge:   int((24 * time.Hour).Seconds()),
		},
	}, nil
}

func logOut(data internal.RouteData[models.LogoutRequest]) (*EmptyRR, error) {
	token := internal.HashString(data.Input.CookieToken)
	if err := data.Database.Where("token = ?", token).Delete(&db.Sessions{}).Error; err != nil {
		slog.Error("error deleting session", "err", err)
		return nil, internal.NewHumaError(http.StatusInternalServerError, "error deleting session")
	}
	db.AddLog(db.LogActionDelete, "User logged out", data.Database)
	return &EmptyRR{}, nil

}
