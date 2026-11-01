package internal

import (
	"log"
	"os"

	"io"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type stdioConn struct {
	closed chan struct{}
	once   sync.Once
}
func (c *stdioConn) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (c *stdioConn) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (c *stdioConn) Close() error {
	c.once.Do(func() { close(c.closed) })
	return nil
}
// LocalAddr, RemoteAddr, Set*Deadline: same stubs as above
type oneConnListener struct {
	conn  *stdioConn
	given bool
}
func (l *oneConnListener) Accept() (net.Conn, error) {
	if !l.given { l.given = true; return l.conn, nil }
	<-l.conn.closed // block until the single connection dies
	return nil, io.EOF
}
func (l *oneConnListener) Close() error   { return nil }
func (l *oneConnListener) Addr() net.Addr { return addr{} }

func Serve() {
	log.SetOutput(os.Stderr) // stdout belongs to the protocol
	srv := grpc.NewServer()
	pb.RegisterAgentServer(srv, &agent{})
	srv.Serve(&oneConnListener{conn: &stdioConn{closed: make(chan struct{})}})
}
