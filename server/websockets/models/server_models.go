package models

type WebSocketServiceMessage string

const (
	Connect          WebSocketServiceMessage = "connect"
	PassPassRequired WebSocketServiceMessage = "pass_pass_required"
)

type AddServiceWSMessage struct {
	Id uint `json:"id" min:"1"`
}
type ServiceWSMessage = WebSocketMessage[WebSocketServiceMessage, any]

func NewServiceWSMessage(action WebSocketServiceMessage, data any) ServiceWSMessage {
	return ServiceWSMessage{Action: action, Data: data}
}
