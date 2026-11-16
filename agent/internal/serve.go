package internal

import (
	"log"
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

func Serve(database *gorm.DB) {
	
	conn := shared.NewStdioConn(os.Stdin, os.Stdout, func() error { return os.Stdin.Close() })
	log.Println("starting agent connection")
	srv := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             15 * time.Second,
		PermitWithoutStream: true,
	}))
	pb.RegisterFilesServiceServer(srv, files.NewFilesServer(database))
	pb.RegisterResourcesServiceServer(srv, &resources.ResourcesServer{})
	pb.RegisterCronServiceServer(srv, cron.NewCronServer(database))

	lis := &oneConnListener{ch: make(chan net.Conn, 1), done: make(chan struct{})}
	log.Println("agent ready")

	lis.ch <- conn // hand the single connection to Accept

	go func() {
		// handle only one connection, stop the server when it's closed
		<-conn.Closed
		srv.Stop()
	}()
	if err := srv.Serve(lis); err != nil {
		log.Println("serve ended:", err)
		log.Println("agent connection closed")
	}
}
