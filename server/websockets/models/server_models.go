package models

type WebSocketServiceMessage string

const (
	ConnectToSSH       WebSocketServiceMessage = "connect_to_ssh"
	PassPhraseRequired WebSocketServiceMessage = "pass_phrase_required"
	Error              WebSocketServiceMessage = "error"
)

type ConnectServerWSMessage struct {
	Id uint `json:"id" min:"1"`
}
type PassPhraseRequiredWSData struct {
	Id   uint `json:"id" min:"1"`
	Name any  `json:"name"`
}
type ServiceWSMessage = WebSocketMessage[WebSocketServiceMessage]

func NewServiceWSMessage(action WebSocketServiceMessage, data any) ServiceWSMessage {
	return ServiceWSMessage{Action: action, Data: data}
}

func NewErrorServiceWSMessage(message string) ServiceWSMessage {
	return NewServiceWSMessage(Error, WSError{Message: message})
}
