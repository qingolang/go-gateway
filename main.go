package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go_gateway/common/lib"
	"go_gateway/dao"
	"go_gateway/grpc_proxy_router"
	"go_gateway/http_proxy_router"
	"go_gateway/router"
	"go_gateway/tcp_proxy_router"
)

var (
	//endpoint dashboard后台管理  server代理服务器
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	//config ./conf/prod/ 对应配置文件夹
	config = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		lib.InitModule(*config)
		defer lib.Destroy()
		router.HTTPServerRun()

		// 接收 Ctrl+c 或 kill 信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HTTPServerStop()
		return
	}
	lib.InitModule(*config)
	defer lib.Destroy()
	// 初始化加载服务列表
	dao.ServiceManagerHandler.LoadOnce()
	// 初始化加载租户
	dao.APPManagerHandler.LoadOnce()

	go func() {
		http_proxy_router.HttpServerRun()
	}()
	go func() {
		http_proxy_router.HttpsServerRun()
	}()
	go func() {
		tcp_proxy_router.TcpServerRun()
	}()
	go func() {
		grpc_proxy_router.GrpcServerRun()
	}()

	// 接收 Ctrl+c 或 kill 信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	tcp_proxy_router.TcpServerStop()
	grpc_proxy_router.GRPCServerStop()
	http_proxy_router.HttpServerStop()
	http_proxy_router.HttpsServerStop()
}
