package http_proxy_middleware

import (
	"go_gateway/common"
	"go_gateway/dao"
	"go_gateway/middleware"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// HTTPJWTAuthTokenMiddleware jwt auth token
func HTTPJWTAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		// 未开启鉴权直接通过
		if serviceDetail.AccessControl.OpenAuth != 1 {
			c.Next()
			return
		}

		// 是否开启API白名单
		if serviceDetail.AccessControl.OpenApiWhiteList == 1 {
			if serviceDetail.AccessControl.ApiWhiteList != "" {
				aptWhiteList := strings.Split(serviceDetail.AccessControl.ApiWhiteList, ",")
				for _, aptWhite := range aptWhiteList {
					if aptWhite == "" {
						continue
					}
					aptWhiteRuneList := []rune(aptWhite)
					if aptWhiteRuneList[len(aptWhiteRuneList)-1] == '*' {
						if strings.HasPrefix(c.Request.URL.Path, strings.TrimSuffix(aptWhite, "*")) {
							c.Next()
							return
						}
					} else {
						if c.Request.URL.Path == aptWhite {
							c.Next()
							return
						}
					}
				}
			}
		}

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := common.JWTDecode(token)
			if err != nil {
				middleware.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			appList := dao.APPManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.APPID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}
		if !appMatched {
			middleware.ResponseError(c, 2003, errors.New("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}
