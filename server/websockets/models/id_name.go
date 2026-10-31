package models

type WebSocketIdNameData struct {
	Id   uint `json:"id" min:"1"`
	Name any  `json:"name"`
}
