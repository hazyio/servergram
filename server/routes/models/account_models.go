package models

import "net/http"

type LoginRequest struct {
	Body struct {
		Password string `json:"password" required:"true" doc:"The password for login" minLength:"16" maxLength:"128"`
	}
	UserAgent string `header:"User-Agent" required:"true" hidden:"true"`
}

type LoginResponse struct {
	Status    int
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type LogoutRequest struct {
	CookieToken string `cookie:"session_token" required:"true" hidden:"true"`
}

type ListSessionsRequest struct {
	CookieToken string `cookie:"session_token" required:"true" hidden:"true"`
}

type Session struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Os        string `json:"os"`
	UserAgent string `json:"user_agent"`
	CreatedAt string  `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	ThisDevice bool   `json:"this_device"`
}

type SessionResponse struct {
	Body struct {
		Sessions []Session `json:"sessions"`
	}
}

type DeleteSessionRequest struct {
	ID string `path:"id" required:"true"`
}
