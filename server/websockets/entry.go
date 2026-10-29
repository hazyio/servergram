package websockets

import (
	"net/http"

	"gorm.io/gorm"
	"hazyio.com/servergram/internal"
)

func RegisterWebSockets(server *http.ServeMux, database *gorm.DB) {
	s := internal.NewRawServer("/ws", database, server)
	s.WebSocket[any, any]("/ws")

}
