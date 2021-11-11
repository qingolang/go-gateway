package http_proxy_middleware

import (
	"go_gateway/dao"
	"go_gateway/middleware"
	"go_gateway/reverse_proxy"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//匹配接入方式 基于请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 加载负载均衡的方式
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		// 获取链接池
		trans, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		// 获取且配置反向代理
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		// 执行
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
