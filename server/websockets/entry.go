package websockets

import (
	"hazyio.com/servergram/internal"
)

func RegisterWebSockets(server *internal.RawServer) {
	server.WebSocket("/services", handleServicesWebSocket)

}
