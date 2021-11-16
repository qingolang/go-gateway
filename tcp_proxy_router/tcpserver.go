package tcp_proxy_router

import (
	"context"
	"fmt"
	"go_gateway/dao"
	"go_gateway/reverse_proxy"
	"go_gateway/tcp_proxy_middleware"
	"go_gateway/tcp_server"
	"log"
)

var tcpServerList = []*tcp_server.TcpServer{}

// TcpServerRun
func TcpServerRun() {
	// 获取服务列表
	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {

			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)
			// 获取加载负载均衡的方式
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [PANIC] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			//构建路由及设置中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				// 流量统计
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				// 限流
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				// 白名单
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				// 黑名单
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)

			//构建回调handler
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}
