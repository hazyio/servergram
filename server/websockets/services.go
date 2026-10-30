package websockets

import (
	"log/slog"

	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/websockets/models"
)

func handleServicesWebSocket(input internal.RawServerRouteData[models.WebSocketMessage[models.WebSocketServiceMessage, any]]) (any, error) {
	message := input.Input.Action
	data := input.Input.Data
	switch message {
	case models.Connect:
		// Handle connect message
		return models.NewServiceWSMessage(models.PassPassRequired, data), nil

	default:
		slog.Error("unknown message action", "action", message)
		return models.NewUnknownError("unknown message action: " + string(message)), nil

	}
	return &data, nil
}
