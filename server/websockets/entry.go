package websockets

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"hazyio.com/servergram/server/internal"
	"hazyio.com/servergram/server/internal/sshmanager"
)

func RegisterWebSockets(api *huma.Group, database *gorm.DB, sshManager *sshmanager.SSHManager, allowedOrigins []string, polls *internal.Polls) {
	server := internal.NewGroup("/ws", api, database, []string{"Services"}, sshManager, allowedOrigins, polls)

	server.WebSocket("/", handleServicesWebSocket, onServiceSocketStated, internal.WebSocketRunAlwaysData{
		Duration: 30 * time.Second,
		Fun:      onServiceSocketStated,
	})

}
