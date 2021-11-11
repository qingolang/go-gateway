package http_proxy_router

import (
	"go_gateway/controller"
	"go_gateway/http_proxy_middleware"
	"go_gateway/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}

	router.Use(
		// 取出当前 service detail
		http_proxy_middleware.HTTPAccessModeMiddleware(),
		// 流量统计
		http_proxy_middleware.HTTPFlowCountMiddleware(),
		// 限流 服务限流 客户端ip限流
		http_proxy_middleware.HTTPFlowLimitMiddleware(),
		// JWT鉴权
		http_proxy_middleware.HTTPJWTAuthTokenMiddleware(),
		// JWT租户流量统计与限流
		http_proxy_middleware.HTTPJWTFlowCountMiddleware(),
		// JWT租户客户端限流
		http_proxy_middleware.HTTPJwtFlowLimitMiddleware(),
		// IP 白名单
		http_proxy_middleware.HTTPWhiteListMiddleware(),
		// IP 黑名单
		http_proxy_middleware.HTTPBlackListMiddleware(),
		// 重写header
		http_proxy_middleware.HTTPHeaderTransferMiddleware(),
		// 去除URL前缀
		http_proxy_middleware.HTTPStripUriMiddleware(),
		// 重写URL
		http_proxy_middleware.HTTPUrlRewriteMiddleware(),
		// 代理
		http_proxy_middleware.HTTPReverseProxyMiddleware())

	return router
}
