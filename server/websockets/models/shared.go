package models

type WebSocketMessage[M, T any] struct {
	Action M `json:"action"`
	Data   T `json:"data"`
}

type WSError struct {
	Error string `json:"error"`
}

func NewUnknownError(err string) WSError {
	return WSError{
		Error: err,
	}
}
