package dto

import (
	"go_gateway/common"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminSessionInfo
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

// AdminLoginInput
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //管理员用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
	Code     string `json:"code" form:"code" comment:"验证码" example:"123456" validate:"required"`                          //验证码
	CodeId   string `json:"codeId" form:"codeId" comment:"验证码ID" example:"123456" validate:"required"`                    //验证码ID
}

// BindValidParam
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return common.DefaultGetValidParams(c, param)
}

// AdminLoginOutput
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` //token
}
