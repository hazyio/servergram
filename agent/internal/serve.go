package internal

import (
	"log"
	"os"

	"net"

	"google.golang.org/grpc"
	"hazyio.com/servergram/agent/internal/files"
	"hazyio.com/servergram/agent/internal/resources"
	"hazyio.com/servergram/shared"
	pb "hazyio.com/servergram/shared/gen/proto"
)

func Serve() {
	log.SetOutput(os.Stderr) // stdout belongs to the protocol
	conn := shared.NewStdioConn(os.Stdin, os.Stdout, func() error { return os.Stdin.Close() })
	log.Println("starting agent connection")
	srv := grpc.NewServer()
	pb.RegisterFilesServer(srv, &files.FilesServer{})
	pb.RegisterServerResourcesServer(srv, &resources.ResourcesServer{})

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
