package models

type WebSocketServiceMessage string

const (
	WSServiceMessageConnectServer           WebSocketServiceMessage = "connect_server"
	WSServiceMessagePassPhraseRequired      WebSocketServiceMessage = "pass_phrase_required"
	WSServiceMessageTrustHostFingerPrint    WebSocketServiceMessage = "trust_host_finger_print"
	WSServiceMessageHostFingerPrintMismatch WebSocketServiceMessage = "host_finger_print_mismatch"
	WSServiceMessageError                   WebSocketServiceMessage = "error"
	WSServiceMessageConnectedServices       WebSocketServiceMessage = "connected_services"
	WSServiceMessageDisconnectServer        WebSocketServiceMessage = "disconnect_server"
)

type WebsocketServiceType string

const (
	WSServiceServerServiceType WebsocketServiceType = "server"
)

type WebsocketServiceStatus string

const (
	WSServiceStatusConnected    WebsocketServiceStatus = "connected"
	WSServiceStatusConnecting   WebsocketServiceStatus = "connecting"
	WSServiceStatusDisconnected WebsocketServiceStatus = "disconnected"
)

type WSServiceConnectServerData struct {
	Id   uint   `json:"id" min:"1"`
	Name string `json:"name"`
}
type WSServiceDisconnectServerData struct {
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
type WSServiceHostFingerPrintMismatchData struct {
	Id             uint   `json:"id" min:"1"`
	Name           any    `json:"name"`
	OldFingerPrint string `json:"old_finger_print"`
	NewFingerPrint string `json:"new_finger_print"`
}
type WSServiceMessage struct {
	Action WebSocketServiceMessage `json:"action" tstype:"'connect_server' | 'pass_phrase_required' | 'trust_host_finger_print' | 'host_finger_print_mismatch' | 'error' | 'connected_services' | 'disconnect_server'"`
	Data   any                     `json:"data"`
}
type WSServiceConnectedListItem struct {
	Id     uint                   `json:"id" min:"1"`
	Name   string                 `json:"name"`
	Type   WebsocketServiceType   `json:"type" tstype:"'server'"`
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
