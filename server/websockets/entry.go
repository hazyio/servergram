package websockets

import (
	"hazyio.com/servergram/server/internal"
)

func RegisterWebSockets(server *internal.Server) {

	server.WebSocket("/services", handleServicesWebSocket,onServiceSocketStated)

}
