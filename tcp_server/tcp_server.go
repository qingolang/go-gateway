package tcp_server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr    string
	Handler TCPHandler
	//err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

func (s *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

func (srv *TcpServer) ListenAndServe() error {
	// 链接是否关闭
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	addr := srv.Addr
	if addr == "" {
		return errors.New("need addr")
	}
	// 获取链接
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Close 关闭
func (srv *TcpServer) Close() error {
	atomic.StoreInt32(&srv.inShutdown, 1)
	close(srv.doneChan) //关闭channel
	srv.l.Close()       //执行listener关闭
	return nil
}

// Serve 服务
func (srv *TcpServer) Serve(l net.Listener) error {
	srv.l = &onceCloseListener{Listener: l}
	defer srv.l.Close() //执行listener关闭
	if srv.BaseCtx == nil {
		srv.BaseCtx = context.Background()
	}

	// 当前server信息写入上下文
	baseCtx := srv.BaseCtx
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			log.Printf(" [ERROR] accept fail, err: %v\n", e)
			continue
		}
		// 获取一个新链接
		c := srv.newConn(rw)
		go c.serve(ctx)
	}
}

// newConn 新链接
func (srv *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

// getDoneChan
func (s *TcpServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}
