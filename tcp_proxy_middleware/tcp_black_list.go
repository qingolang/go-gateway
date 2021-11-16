package tcp_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"strings"
)

// TCPBlackListMiddleware 匹配接入方式 基于请求信息
func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		if serviceDetail.AccessControl.OpenBlackList != 1 {
			c.Next()
			return
		}
		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		// 当前IP如果存在于白名单则不校验黑名单
		if serviceDetail.AccessControl.WhiteList != "" {
			for _, ipWhite := range strings.Split(serviceDetail.AccessControl.WhiteList, ",") {
				if clientIP == ipWhite {
					c.Next()
					return
				}
			}
		}
		// 校验黑名单
		if serviceDetail.AccessControl.BlackList != "" {
			if common.InStringSlice(strings.Split(serviceDetail.AccessControl.BlackList, ","), clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
