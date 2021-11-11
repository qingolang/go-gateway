package http_proxy_middleware

import (
	"go_gateway/dao"
	"go_gateway/middleware"

	"github.com/gin-gonic/gin"
)

// HTTPAccessModeMiddleware 匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		c.Set("service", service)
		c.Next()
	}
}
