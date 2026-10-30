package models

type ServerMessage string

const (
	Connect ServerMessage = "connect"
)

type ServiceWSMessages[T any] struct {
	Message ServerMessage `json:"message"`
	Data    T             `json:"data"`
}

type AddServiceWSMessage struct {
	Id uint `json:"id" min:"1"`
}
