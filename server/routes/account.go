package routes

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/internal"
)

func RegisterAccountRoutes(api *huma.Group, database *gorm.DB, loginTokens *internal.LoginToken) {
	server := internal.NewGroup("/account", api, database)
	server.PostRaw("/login", http.StatusSeeOther, func(nrd internal.NormalRouteData[internal.BodyRequest[LoginRequestBody]]) (*LoginSuccessfulRedirectResponse, error) {
		return login(nrd.Input.Body, loginTokens)
	})

}

func login(data LoginRequestBody, loginTokens *internal.LoginToken) (*LoginSuccessfulRedirectResponse, error) {
	env := internal.LoadEnv()
	if env.Password == "" {
		slog.Error("server password is not configured")
		return nil, internal.NewHumaError(http.StatusInternalServerError, "server password is not configured")
	}
	a := sha256.Sum256([]byte(data.Password))
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

	loginTokens.Add(randomToken, expires)

	return &LoginSuccessfulRedirectResponse{
		Status:   http.StatusSeeOther,
		Location: "/",
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
