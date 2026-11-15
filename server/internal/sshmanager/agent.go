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
	return &Agent{Conn: conn, Timeout: 10 * time.Second}
}


func (c *Agent) Close() error {
	if c == nil || c.Conn == nil {
		return nil
	}
	return c.Conn.Close()
}

func (c *Agent) UpdateTimeout(seconds uint8) {
	c.Timeout = time.Duration(seconds) * time.Second
}
func (c *Agent) ResourcesClient(ctx context.Context) (proto.ResourcesServiceClient, context.Context, context.CancelFunc) {
	slog.Debug("Creating server resources client")
	ctx, cancel := context.WithCancel(ctx)
	return proto.NewResourcesServiceClient(c.Conn), ctx, cancel
}

func (c *Agent) FilesClient(ctx context.Context) (proto.FilesServiceClient, context.Context, context.CancelFunc) {
	slog.Debug("Creating files client")
	ctx, cancel := context.WithCancel(ctx)
	return proto.NewFilesServiceClient(c.Conn), ctx, cancel
}

func (c *Agent) CronClient(ctx context.Context) (proto.CronServiceClient, context.Context, context.CancelFunc) {
	slog.Debug("Creating cron client")
	ctx, cancel := context.WithCancel(ctx)
	return proto.NewCronServiceClient(c.Conn), ctx, cancel
}
