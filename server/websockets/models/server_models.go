package models

type ServerMessage string

const (
	Add ServerMessage = "add"
)

type ServerWebSocketMessages struct {
	Message ServerMessage `json:"message"`
}
