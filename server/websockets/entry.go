package websockets

import (
	"context"

	"hazyio.com/servergram/internal"
)

func RegisterWebSockets(server *internal.RawServer) {
	server.WebSocket("/services", handleServicesWebSocket)

}
