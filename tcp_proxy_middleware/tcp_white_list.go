package tcp_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"strings"
)

// TCPWhiteListMiddleware 匹配接入方式 基于请求信息
func TCPWhiteListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.OpenWhiteList != 1 {
			c.Next()
			return
		}

		// 取出IP
		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		// 校验白名单
		if serviceDetail.AccessControl.WhiteList != "" {
			if !common.InStringSlice(strings.Split(serviceDetail.AccessControl.WhiteList, ","), clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
