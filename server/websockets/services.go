package websockets

import (
	"log/slog"

	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/websockets/models"
)

func handleServicesWebSocket(input internal.RawServerRouteData[models.WebSocketMessage[models.WebSocketServiceMessage]]) (any, error) {
	message := input.Input.Action
	data := input.Input.Data
	switch message {
	case models.Connect:
		msg := models.ConnectServerWSMessage{}
		if err := models.ReUnmarshal(data, &msg); err != nil {
			return models.NewErrorServiceWSMessage("could not understand message: " + string(message)), nil
		}
		return handleConnect(msg)

	default:
		slog.Error("unknown message action", "action", message)
		return models.NewErrorServiceWSMessage("unknown message action: " + string(message)), nil

	}
}

func handleConnect(msg models.ConnectServerWSMessage) (any, error) {
	// marshaledData, err := models.Marshal(&models.PassPhraseRequiredWSData{
	// 	Id:   msg.Id,
	// 	Name: "Service Name",
	// })
	// if err != nil {
	// 	return models.NewErrorServiceWSMessage("could not marshal response data"), nil
	// }
	return models.NewServiceWSMessage(models.PassPhraseRequired, &models.PassPhraseRequiredWSData{
		Id:   msg.Id,
		Name: "Service Name",
	}), nil
}
