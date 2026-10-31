package models

type WebSocketServiceMessage string

const (
	WSServiceMessageConnectServer        WebSocketServiceMessage = "connect_server"
	WSServiceMessagePassPhraseRequired   WebSocketServiceMessage = "pass_phrase_required"
	WSServiceMessageTrustHostFingerPrint WebSocketServiceMessage = "trust_host_finger_print"
	WSServiceMessageError                WebSocketServiceMessage = "error"
	WSServiceMessageConnectedServices    WebSocketServiceMessage = "connected_services"
)

type WebsocketServiceType string

const (
	WSServiceServerServiceType WebsocketServiceType = "server"
)

type WebsocketServiceStatus string

const (
	WSServiceStatusConnected    WebsocketServiceStatus = "connected"
	WSServiceStatusPending      WebsocketServiceStatus = "pending"
	WSServiceStatusConnecting   WebsocketServiceStatus = "connecting"
	WSServiceStatusDisconnected WebsocketServiceStatus = "disconnected"
)

type WSServiceConnectServerData struct {
	Id uint `json:"id" min:"1"`
}
type WSServicePassPhraseRequiredData struct {
	Id   uint `json:"id" min:"1"`
	Name any  `json:"name"`
}
type WSServiceTrustHostFingerPrintData struct {
	Id          uint   `json:"id" min:"1"`
	Name        any    `json:"name"`
	FingerPrint string `json:"finger_print"`
}
type WSServiceMessage = WebSocketMessage[WebSocketServiceMessage]
type WSServiceConnectedListItem struct {
	Id     uint                   `json:"id" min:"1"`
	Name   string                 `json:"name"`
	Type   WebsocketServiceType   `json:"type"`
	Status WebsocketServiceStatus `json:"status" tstype:"'connected' | 'pending' | 'connecting' | 'disconnected'"`
}

type WSServiceConnectedList struct {
	Items []WSServiceConnectedListItem `json:"items"`
}

func NewServiceWSMessage(action WebSocketServiceMessage, data any) WSServiceMessage {
	return WSServiceMessage{Action: action, Data: data}
}

func NewErrorServiceWSMessage(message string) WSServiceMessage {
	return NewServiceWSMessage(WSServiceMessageError, WSError{Message: message})
}
