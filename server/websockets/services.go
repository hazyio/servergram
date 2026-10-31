package websockets

import (
	"log/slog"

	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/websockets/models"
)

func handleServicesWebSocket(input internal.RawServerRouteData[models.WebSocketMessage[models.WebSocketServiceMessage]]) {
	message := input.Input.Action
	data := input.Input.Data
	dataQueue := input.DataQueue
	switch message {
	case models.ConnectToSSH:
		msg := models.ConnectServerWSMessage{}
		if err := models.ReUnmarshal(data, &msg); err != nil {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("could not understand message: " + string(message)))
			// return models.NewErrorServiceWSMessage("could not understand message: " + string(message)), nil
		}
		handleConnect(msg, dataQueue)

	default:
		slog.Error("unknown message action", "action", message)
		dataQueue.Enqueue(models.NewErrorServiceWSMessage("unknown message action: " + string(message)))
		// return models.NewErrorServiceWSMessage("unknown message action: " + string(message)), nil

	}
}

func handleConnect(msg models.ConnectServerWSMessage, dataQueue *internal.DataQueue) {
	dataQueue.Enqueue(models.NewServiceWSMessage(models.PassPhraseRequired, &models.PassPhraseRequiredWSData{
		Id:   msg.Id,
		Name: "Service Name",
	}))
}
