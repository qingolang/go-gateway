package controller

import (
	"go_gateway/common/captcha"
	"go_gateway/dto"
	"go_gateway/middleware"

	"github.com/gin-gonic/gin"
)

// CaptchaController
type CaptchaController struct{}

// CaptchaRegister
func CaptchaRegister(group *gin.RouterGroup) {
	Captacha := &CaptchaController{}
	group.GET("/get", Captacha.Captacha)
}

// Captacha godoc
// @Summary 获取验证码
// @Description 获取验证码
// @Tags 验证码接口
// @ID /captacha/
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.CaptchaOutput} "success"
// @Router /captacha/ [get]
func (Captacha *CaptchaController) Captacha(c *gin.Context) {
	id, image, err := captcha.GenerateCaptchaHandler()
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	out := &dto.CaptchaOutput{ID: id, Image: image}
	middleware.ResponseSuccess(c, out)
}
