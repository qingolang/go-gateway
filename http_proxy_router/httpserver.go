package http_proxy_router

import (
	"context"
	"go_gateway/common/lib"
	"go_gateway/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

// HttpServerRun
func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf("[INFO]  http_proxy_run %s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR]  http_proxy_run %s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}
}

// HttpsServerRun
func HttpsServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf("[INFO] https_proxy_run %s\n", lib.GetStringConf("proxy.https.addr"))
	// todo 以下命令只在编译机有效，如果是交叉编译情况下需要单独设置路径
	//if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err!=http.ErrServerClosed {
	if err := HttpsSrvHandler.ListenAndServeTLS(lib.GetStringConf("proxy.https.cert_file_crt"),
		lib.GetStringConf("proxy.https.cert_file_key")); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] https_proxy_run %s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}
}

// HttpServerStop
func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] http_proxy_stop err:%v\n", err)
	}
	cancel()
	log.Printf("[INFO] http_proxy_stop %v stopped\n", lib.GetStringConf("proxy.http.addr"))
}

// HttpsServerStop
func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	cancel()
	log.Printf("[INFO] https_proxy_stop %v stopped\n", lib.GetStringConf("proxy.https.addr"))
}
