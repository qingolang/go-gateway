package http_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"go_gateway/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HTTPWhiteListMiddleware 匹配接入方式 基于请求信息
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		// 是否开启校验
		if serviceDetail.AccessControl.OpenWhiteList != 1 {
			c.Next()
			return
		}

		// 校验白名单
		if serviceDetail.AccessControl.WhiteList != "" {
			if !common.InStringSlice(strings.Split(serviceDetail.AccessControl.WhiteList, ","), c.ClientIP()) {
				middleware.ResponseError(c, 3001, errors.New(fmt.Sprintf("%s not in white ip list", c.ClientIP())))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
