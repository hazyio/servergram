package models

import "net/http"

type LoginRequest struct {
	Body struct {
		Password string `json:"password" required:"true" doc:"The password for login" minLength:"16" maxLength:"128"`
	}
	UserAgent string `header:"User-Agent" required:"true" hidden:"true"`
}

type LoginSuccessfulRedirectResponse struct {
	Status    int
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type LogoutRequest struct {
	CookieToken string `header:"Cookie-Token" required:"true" hidden:"true"`
}

