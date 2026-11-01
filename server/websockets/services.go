package websockets

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
	"hazyio.com/servergram/server/db"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
	"hazyio.com/servergram/server/websockets/models"
)

func onServiceSocketStated(input internal.OnStartedRawServerRouteData) {
	manageSSHConnection := input.SSHManager
	dataQueue := input.DataQueue
	SendQueueAndConnected(dataQueue, manageSSHConnection)

}
func handleServicesWebSocket(input internal.RawServerRouteData[models.WebSocketMessage[models.WebSocketServiceMessage]]) {
	message := input.Input.Action
	data := input.Input.Data
	manageSSHConnection := input.SSHManager
	dataQueue := input.DataQueue
	slog.Info("handling service websocket message", "action", message)
	switch message {
	case models.WSServiceMessageConnectServer:
		msg := models.WSServiceConnectServerData{}
		if err := models.ReUnmarshal(data, &msg); err != nil {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("could not understand message: " + string(message)))
		}
		if manageSSHConnection.IsQueued(msg.Id) {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("server is already queued"))
			return
		}
		if manageSSHConnection.IsConnected(msg.Id) {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("server is already connected"))
			return
		}
		manageSSHConnection.Queue(msg.Id, msg.Name)
		db.AddLog(db.Update, "Queued connection to server with ID: "+fmt.Sprint(msg.Id), input.Database)
		SendQueueAndConnected(dataQueue, manageSSHConnection)
		// Handle the server connection asynchronously.
		go handleServerConnect(manageSSHConnection, input.Database, msg, dataQueue)

	case models.WSServiceMessageDisconnectServer:
		msg := models.WSServiceDisconnectServerData{}
		if err := models.ReUnmarshal(data, &msg); err != nil {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("could not understand message: " + string(message)))
			return
		}
		go handleServerDisconnect(manageSSHConnection, msg, dataQueue)

	default:
		slog.Error("unknown message action", "action", message)
		dataQueue.Enqueue(models.NewErrorServiceWSMessage("unknown message action: " + string(message)))

	}
}
func handleServerDisconnect(sshManager *sshmanager.SSHManager, msg models.WSServiceDisconnectServerData, dataQueue *internal.DataQueue) {
	defer SendQueueAndConnected(dataQueue, sshManager)

	if sshManager.IsConnected(msg.Id) {
		if err := sshManager.Disconnect(msg.Id); err != nil {
			slog.Error("failed to disconnect server", "error", err)
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("failed to disconnect server"))
			return
		}
	}
}
func handleServerConnect(sshManager *sshmanager.SSHManager, database *gorm.DB, msg models.WSServiceConnectServerData, dataQueue *internal.DataQueue) {
	var server db.Server
	defer SendQueueAndConnected(dataQueue, sshManager)
	defer sshManager.Dequeue(msg.Id)

	if err := database.Preload("SSHKey").Where("id = ?", msg.Id).First(&server).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dataQueue.Enqueue(models.NewErrorServiceWSMessage("could not find server"))
			return
		}
		slog.Error("unexpected error occurred while finding server", "error", err)
		dataQueue.Enqueue(models.NewErrorServiceWSMessage("An unexpected error occurred"))
		return
	}
	if server.SSHKey.Encrypted {
		dataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageTrustHostFingerPrint, &models.WSServicePassPhraseRequiredData{
			Id:   msg.Id,
			Name: server.Name,
		}))
		return
	}

	env := internal.LoadEnv()
	decryptedSSHKey, decryptErr := internal.DecryptString(server.SSHKey.Key, env.EncryptionMasterKey)
	if decryptErr != nil {
		slog.Error("failed to decrypt SSH key", "error", decryptErr)
		dataQueue.Enqueue(models.NewErrorServiceWSMessage("failed to decrypt SSH key"))

		return
	}
	server.SSHKey.Key = decryptedSSHKey
	connectErr := sshManager.ConnectNewUnEncrypted(&server)
	if connectErr == nil {
		// get and update os and arch information
		client := sshManager.GetConnected(msg.Id)
		if client != nil {
			os := client.OS
			arch := client.Arch
			if err := db.UpdateServerOsAndArch(msg.Id, os, arch, database); err != nil {
				slog.Error("failed to update server OS and Arch", "error", err)
			}
		}

		return
	}
	var unknownHost *sshmanager.UnknownHostError
	var mismatch *sshmanager.HostKeyMismatchError
	if errors.As(connectErr, &unknownHost) {
		slog.Error("unknown host error", "error", unknownHost)
		dataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageTrustHostFingerPrint, &models.WSServiceTrustHostFingerPrintData{
			Id:          msg.Id,
			Name:        server.Name,
			FingerPrint: unknownHost.Fingerprint,
		}))

		return
	}
	if errors.As(connectErr, &mismatch) {
		slog.Error("host key mismatch error", "error", mismatch)
		dataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageHostFingerPrintMismatch, &models.WSServiceHostFingerPrintMismatchData{
			Id:             msg.Id,
			Name:           server.Name,
			OldFingerPrint: mismatch.Expected,
			NewFingerPrint: mismatch.Got,
		}))

		return
	}
	slog.Error("failed to connect to server", "error", connectErr)
	dataQueue.Enqueue(models.NewErrorServiceWSMessage("failed to connect to server, config might be wrong. check server configuration"))

}

func SendQueueAndConnected(dataQueue *internal.DataQueue, sshManager *sshmanager.SSHManager) {
	connectedQueue, connectedServices := sshManager.GetQueue(), sshManager.GetAllConnected()
	data := make([]models.WSServiceConnectedListItem, 0)
	// Map connected services to WSServiceConnectedListItem and append to data
	data = append(data, internal.Map(connectedServices, func(item sshmanager.SSHConnectedListData) models.WSServiceConnectedListItem {
		return models.WSServiceConnectedListItem{
			Id:     item.Id,
			Name:   item.Name,
			Type:   models.WSServiceServerServiceType,
			Status: models.WSServiceStatusConnected,
		}
	})...)
	// Map queued services to WSServiceConnectedListItem
	data = append(data, internal.Map(connectedQueue, func(item sshmanager.SSHConnectedListData) models.WSServiceConnectedListItem {
		return models.WSServiceConnectedListItem{
			Id:     item.Id,
			Name:   item.Name,
			Type:   models.WSServiceServerServiceType,
			Status: models.WSServiceStatusConnecting,
		}
	})...)

	// Enqueue the list of connected and queued services to the data queue
	dataQueue.Enqueue(models.NewServiceWSMessage(models.WSServiceMessageConnectedServices, &models.WSServiceConnectedList{
		Items: data,
	}))
}
