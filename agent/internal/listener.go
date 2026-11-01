package internal

import (
	"net"
	"sync"

	"hazyio.com/servergram/shared"
)
type oneConnListener struct {
    ch   chan net.Conn
    done chan struct{}
    once sync.Once
}

func (l *oneConnListener) Accept() (net.Conn, error) {
    select {
    case c := <-l.ch:
        return c, nil
    case <-l.done:
        return nil, net.ErrClosed
    }
}
func (l *oneConnListener) Close() error { l.once.Do(func() { close(l.done) }); return nil }
func (l *oneConnListener) Addr() net.Addr { return shared.StdIoAddr{} }
