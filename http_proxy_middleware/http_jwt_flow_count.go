package http_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"go_gateway/middleware"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HTTPJWTFlowCountMiddleware JWT租户流量统计 与 限流操作
func HTTPJWTFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.APP)
		appCounter, err := common.FlowCounterHandler.GetCounter(common.FlowAppPrefix + appInfo.APPID)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.QPD, appCounter.TotalCount)))
			c.Abort()
			return
		}
		c.Next()
	}
}
