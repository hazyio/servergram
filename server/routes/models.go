package routes

import "net/http"

type LoginRequestBody struct {
	Body struct {
		Password string `json:"password" required:"true" doc:"The password for login" minLength:"16" maxLength:"128"`
	}
	UserAgent string `header:"User-Agent" required:"true"`
}

type LoginSuccessfulRedirectResponse struct {
	Status    int
	SetCookie http.Cookie `header:"Set-Cookie"`
}
type MessageResponse struct {
	Message string `json:"message"`
}
