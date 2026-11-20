package internal

import (
	"io"
	"log/slog"
	"os"
	"time"

	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/cron"
	"hazyio.com/servergram/agent/internal/files"
	"hazyio.com/servergram/agent/internal/resources"
	"hazyio.com/servergram/shared"
	pb "hazyio.com/servergram/shared/gen/proto"
)

func Serve(database *gorm.DB, transportOut io.Writer) {
	conn := shared.NewStdioConn(os.Stdin, transportOut, func() error { return os.Stdin.Close() })
	slog.Info("starting agent connection")
	srv := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             15 * time.Second,
		PermitWithoutStream: true,
	}))
	pb.RegisterFilesServiceServer(srv, files.NewFilesServer(database))
	pb.RegisterResourcesServiceServer(srv, &resources.ResourcesServer{})
	pb.RegisterCronServiceServer(srv, cron.NewCronServer(database))

	lis := &oneConnListener{ch: make(chan net.Conn, 1), done: make(chan struct{})}
	slog.Info("agent ready")

	lis.ch <- conn // hand the single connection to Accept

	go func() {
		// handle only one connection, stop the server when it's closed
		<-conn.Closed
		srv.Stop()
	}()
	if err := srv.Serve(lis); err != nil {
		slog.Info("serve ended", "error", err)
		slog.Info("agent connection closed")
	}
}
