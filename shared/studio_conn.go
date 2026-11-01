package shared

import (
	"io"
	"net"
	"sync"
	"time"
)

type StdIoAddr struct{}

func (StdIoAddr) Network() string { return "stdio" }
func (StdIoAddr) String() string  { return "stdio" }

type StdIoConn struct {
	io.Reader
	io.Writer
	closeFn func() error
	once    sync.Once
	Closed  chan struct{}
}

func NewStdioConn(r io.Reader, w io.Writer, closeFn func() error) *StdIoConn {
	return &StdIoConn{Reader: r, Writer: w, closeFn: closeFn, Closed: make(chan struct{})}
}

func (c *StdIoConn) Close() error {
	var err error
	c.once.Do(func() { close(c.Closed); err = c.closeFn() })
	return err
}
func (*StdIoConn) LocalAddr() net.Addr              { return StdIoAddr{} }
func (*StdIoConn) RemoteAddr() net.Addr             { return StdIoAddr{} }
func (*StdIoConn) SetDeadline(time.Time) error      { return nil }
func (*StdIoConn) SetReadDeadline(time.Time) error  { return nil }
func (*StdIoConn) SetWriteDeadline(time.Time) error { return nil }
