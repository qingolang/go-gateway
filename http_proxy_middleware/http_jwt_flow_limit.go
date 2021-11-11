package http_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"go_gateway/middleware"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HTTPJwtFlowLimitMiddleware
func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.APP)
		if appInfo.QPS > 0 {
			clientLimiter, err := common.FlowLimiterHandler.GetLimiter(
				common.FlowAppPrefix+appInfo.APPID+"_"+c.ClientIP(),
				float64(appInfo.QPS))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				middleware.ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.QPS)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
