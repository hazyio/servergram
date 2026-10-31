package websockets

import (
	"log/slog"

	"hazyio.com/servergram/internal"
	"hazyio.com/servergram/websockets/models"
)

func onServiceSocketStated(input internal.OnStartedRawServerRouteData) {
	services := &models.WSServiceConnectedList{}
	services.Items = append(services.Items, models.WSServiceConnectedListItem{
		Id:     7,
		Name:   "Service Name",
		Type:   models.WSServiceServerServiceType,
		Status: models.WSServiceStatusConnecting,
	})
	// Send the connected services list to the data queue
	input.DataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageConnectedServices, services))
}
func handleServicesWebSocket(input internal.RawServerRouteData[models.WebSocketMessage[models.WebSocketServiceMessage]]) {
	message := input.Input.Action
	data := input.Input.Data
	dataQueue := input.DataQueue
	switch message {
	case models.WSServiceMessageConnectServer:
		msg := models.WSServiceConnectServerData{}
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

func handleConnect(msg models.WSServiceConnectServerData, dataQueue *internal.DataQueue) {

	dataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageTrustHostFingerPrint, &models.WSServiceTrustHostFingerPrintData{
		Id:   msg.Id,
		Name: "Service Name",
		FingerPrint: `|1|JT3PpY2BaEkKeNC/fWlR3sNPZEI=|vOdGmBIJrpTTbQtHB4069KAm+8I= ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGW2qFypeAHYwUybyrEQ8xLM5/1rMkJK3Eaz1kVd9EA6
|1|M4X9cDDa7n9k00xrW6fl3T5VSP0=|b/hNbhVT/HAc/Cggl24SfIsg+gM= ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDEsavoTB5IBeSdVqn25+Ps2fvHXTBhv5C/J+6f40wgudoAeRbPUBfLaxFFh8YGG/8x+u4WitvrkCp3o5nw4euhF7gJzLfsI3z7ZC2L+sxUsZ8s2lzNRpdbXJchd1q0Jt2cq948uqZrHIDU5aHJMkK0chr6XNzheFDoy9OZ9aBihI54RVQtqamvXbOSIkXan68pxK24g2RUXWXu8a4h1vzmC/AuRQuWsVy694qzyqi3fzLxrIH3D2+BtFb8cMQngQhArFBM7jPokT0XPcaYhrgCUyoOuzHeYQaQzNfG8+xyDRDJFSG6Rvcy3FmczQCOFTPFHAKCnl2672wA5acymvt6lXg/VcD8aVtI+vsQVy0PGi6BfiGtNd9moEg6GAa779hPlDeCaZhyvXJOmW/F2hsSA/Ec9jUdHPj/Tvr8HHLFDfxpZwd6wD72Um1Ve/T69Kzinf4Jsg3HdiGQ3xnB4eMaOgn8AUVI56oMYb4vGxbxvWGX9+oCOJONvHnJ9OiQdw8=
|1|1BSJG2+HYXA435yV1u/tdCpoErs=|8e9y/OsyQBPEIcJq7yOG3d6YiDE= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH1kpe3AXocksgEn9NvI+fzruRluz0EI1ZIabTRQL03+mCOzSIh0zIDRBM14biOGQv7o4hmVIsTAuXicF84OAXk=
`,
	}))
}
