package websockets

import (
	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/websockets/models"
)

func handleServicesWebSocket[I any](input internal.RawServerRouteData[models.ServiceWSMessages[any]]) (any, error) {
	message := input.Input.Message
	data := input.Input.Data
	switch message {
	case models.Connect:
		// Handle connect message
	}
	return &data, nil
}
