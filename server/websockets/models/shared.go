package models

import (
	"encoding/json"
)

type WebSocketMessage[A any] struct {
	Action A   `json:"action"`
	Data   any `json:"data"`
}

type WSError struct {
	Message string `json:"message"`
}
type WebSocketIdNameData struct {
	Id   uint `json:"id" min:"1"`
	Name any  `json:"name"`
}

func NewUnknownError(err string) WSError {
	return WSError{
		Message: err,
	}
}

func ReUnmarshal[T any](data any, v *T) error {
	// very wasteful, but necessary to convert from any to the desired type
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func Marshal[T any](v *T) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
