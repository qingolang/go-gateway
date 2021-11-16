package controller

import (
	"encoding/json"
	"errors"
	"time"

	"go_gateway/common"
	"go_gateway/common/captcha"
	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AdminLoginController
type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {

	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 验证码校验
	if !captcha.CaptchaVerifyHandle(params.CodeId, params.Code) {
		middleware.ResponseError(c, 2001, errors.New("验证码错误"))
		return
	}

	//1. params.UserName 取得管理员信息 admininfo
	//2. admininfo.salt + params.Password sha256 => saltPassword
	//3. saltPassword==admininfo.password
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(common.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(common.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
