package router

import (
	"go_gateway/common/lib"
	"go_gateway/controller"
	"go_gateway/docs"
	"go_gateway/middleware"
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
// InitRouter
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// swagger Init
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	// ??????
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// swagger API ??????????????????
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ??????session
	store, err := sessions.NewRedisStore(10, "tcp", lib.GetStringConf("base.session.redis_server"), lib.GetStringConf("base.session.redis_password"), []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}
	// ???????????????
	adminLoginRouter := router.Group("/admin_login")
	adminLoginRouter.Use(
		// session
		sessions.Sessions("mysession", store),
		// ??????panic
		middleware.RecoveryMiddleware(),
		// ??????????????????
		middleware.RequestLog(),
		// ?????????
		middleware.TranslationMiddleware())
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}

	// ?????????
	captchaRouter := router.Group("/captcha")
	captchaRouter.Use(
		// session
		sessions.Sessions("mysession", store),
		// ??????panic
		middleware.RecoveryMiddleware(),
		// ??????????????????
		middleware.RequestLog(),
		// ?????????
		middleware.TranslationMiddleware())
	{
		controller.CaptchaRegister(captchaRouter)
	}

	// ???????????????
	adminRouter := router.Group("/admin")
	adminRouter.Use(
		// session
		sessions.Sessions("mysession", store),
		// ??????panic
		middleware.RecoveryMiddleware(),
		// ????????????
		middleware.RequestLog(),
		// session ??????
		middleware.SessionAuthMiddleware(),
		// ??????
		middleware.TranslationMiddleware())
	{
		controller.AdminRegister(adminRouter)
	}

	// ????????????
	serviceRouter := router.Group("/service")
	serviceRouter.Use(
		// ??????session
		sessions.Sessions("mysession", store),
		// ??????panic
		middleware.RecoveryMiddleware(),
		// ????????????
		middleware.RequestLog(),
		// session ??????
		middleware.SessionAuthMiddleware(),
		// ??????
		middleware.TranslationMiddleware())
	{
		controller.ServiceRegister(serviceRouter)
	}

	// ????????????
	appRouter := router.Group("/app")
	appRouter.Use(
		// ??????session
		sessions.Sessions("mysession", store),
		// ??????panic
		middleware.RecoveryMiddleware(),
		// ????????????
		middleware.RequestLog(),
		// session ??????
		middleware.SessionAuthMiddleware(),
		// ??????
		middleware.TranslationMiddleware())
	{
		controller.APPRegister(appRouter)
	}

	// ??????
	dashRouter := router.Group("/dashboard")
	dashRouter.Use(
		// ?????? Session
		sessions.Sessions("mysession", store),
		// ?????? panic
		middleware.RecoveryMiddleware(),
		// ????????????
		middleware.RequestLog(),
		// session ??????
		middleware.SessionAuthMiddleware(),
		// ??????
		middleware.TranslationMiddleware())
	{
		controller.DashboardRegister(dashRouter)
	}
	// router.Static("/dist", "./dist")
	return router
}
