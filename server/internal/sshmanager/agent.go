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
	return &Agent{Conn: conn, Timeout: 30 * time.Second}
}

func (c *Agent) UpdateTimeout(seconds uint8) {
	c.Timeout = time.Duration(seconds) * time.Second
}
func (c *Agent) ServerResourcesClient(ctx context.Context) (proto.ServerResourcesClient, context.Context, context.CancelFunc) {
	slog.Debug("Creating server resources client")
	ctx, cancel := context.WithCancel(ctx)
	return proto.NewServerResourcesClient(c.Conn), ctx, cancel
}

func (c *Agent) FilesClient(ctx context.Context) (proto.FilesClient, context.Context, context.CancelFunc) {
	slog.Debug("Creating files client")
	ctx, cancel := context.WithCancel(ctx)
	return proto.NewFilesClient(c.Conn), ctx, cancel
}
