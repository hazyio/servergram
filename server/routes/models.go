package routes

import "net/http"


type LoginRequestBody struct {
	Password string `json:"password" required:"true" doc:"The password for login" minLength:"16" maxLength:"128"`
}

type LoginSuccessfulRedirectResponse struct {
	Status    int
	Location  string      `header:"Location"`
	SetCookie http.Cookie `header:"Set-Cookie"`
}
type MessageResponse struct {
	Message string `json:"message"`
}
