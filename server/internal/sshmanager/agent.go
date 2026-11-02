package sshmanager

import (
	"log/slog"
	"time"

	"context"

	"google.golang.org/grpc"
	"hazyio.com/servergram/shared/gen/proto"
)

type Agent struct {
	Conn    *grpc.ClientConn
	Timeout time.Duration
}

// func (c *Agent) Close() {

// }
func NewAgent(conn *grpc.ClientConn) *Agent {
	return &Agent{Conn: conn, Timeout: time.Duration(10) * time.Second}
}

func (c *Agent) UpdateTimeout(seconds uint8) {
	c.Timeout = time.Duration(seconds) * time.Second
}
func (c *Agent) ServerResourcesClient(ctx context.Context) (proto.ServerResourcesClient, context.Context, context.CancelFunc) {
	slog.Info("Creating server resources client")
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	return proto.NewServerResourcesClient(c.Conn), ctx, cancel
}
