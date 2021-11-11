package grpc_proxy_router

import (
	"fmt"
	"go_gateway/dao"
	"go_gateway/grpc_proxy_middleware"
	"go_gateway/reverse_proxy"
	"log"
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

// grpcServerList
var grpcServerList = []*warpGRPCServer{}

// warpGRPCServer
type warpGRPCServer struct {
	Addr string
	*grpc.Server
}

// GrpcServerRun
func GrpcServerRun() {
	// 获取服务列表
	serviceList := dao.ServiceManagerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			// 获取负载均衡策略
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [INFO] GrpcListen %v err:%v\n", addr, err)
			}
			// 构建
			grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				// grpc 服务拦截器
				grpc.ChainStreamInterceptor(
					// 流量统计
					grpc_proxy_middleware.GRPCFlowCountMiddleware(serviceDetail),
					// 限流 服务限流 客户端限流
					grpc_proxy_middleware.GRPCFlowLimitMiddleware(serviceDetail),
					// JWT 鉴权
					grpc_proxy_middleware.GRPCJWTAuthTokenMiddleware(serviceDetail),
					// JWT 租户流量统计与租户限流
					grpc_proxy_middleware.GRPCJWTFlowCountMiddleware(serviceDetail),
					// JWT 租户 客户端限流
					grpc_proxy_middleware.GRPCJWTFlowLimitMiddleware(serviceDetail),
					// 白名单
					grpc_proxy_middleware.GRPCWhiteListMiddleware(serviceDetail),
					// 黑名单
					grpc_proxy_middleware.GRPCBlackListMiddleware(serviceDetail),
					// 重写header
					grpc_proxy_middleware.GRPCHeaderTransferMiddleware(serviceDetail),
				),
				// 获取消息编码与解码器
				grpc.CustomCodec(proxy.Codec()),
				// 处理服务器流包括上下文
				grpc.UnknownServiceHandler(grpcHandler))

			grpcServerList = append(grpcServerList, &warpGRPCServer{
				Addr:   addr,
				Server: s,
			})
			log.Printf("[INFO] grpc_proxy_run %v\n", addr)
			if err := s.Serve(lis); err != nil {
				log.Fatalf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

// GRPCServerStop
func GRPCServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf("[INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
