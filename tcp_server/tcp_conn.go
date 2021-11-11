package tcp_server

import (
	"context"
	"log"
	"net"
	"runtime"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tcp_proxy context value " + k.name
}

type conn struct {
	server *TcpServer
	//cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) close() {
	c.rwc.Close()
}

// serve 进一步处理
func (c *conn) serve(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10 // 1024
			buf := make([]byte, size)
			// 获取运行日志
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("[PANIC] panic serving %v: %v\n%s\n", c.remoteAddr, err, buf)
		}
		c.close()
	}()
	// 将当前连接地址写入上下文
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	if c.server.Handler == nil {
		panic("handler empty")
	}
	// 回调 routerHandler
	c.server.Handler.ServeTCP(ctx, c.rwc)
}
