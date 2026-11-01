package websockets

import (
	"hazyio.com/servergram/server/internal"
)

func RegisterWebSockets(server *internal.RawServer) {

	server.WebSocket("/services", handleServicesWebSocket,onServiceSocketStated)

}
