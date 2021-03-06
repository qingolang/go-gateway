package middleware

import (
	"errors"
	"fmt"
	"go_gateway/common"
	"go_gateway/common/lib"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.ComLogNotice(c, "_com_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})
				if lib.ConfBase.DebugMode != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				} else {
					ResponseError(c, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
