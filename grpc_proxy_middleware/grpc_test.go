package grpc_proxy_middleware

import (
	"log"

	"google.golang.org/grpc"
)

// GRPCAuthStreamInterceptor 流式RPC拦截器
func GRPCAuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	//md, ok := metadata.FromIncomingContext(ss.Context())
	err := handler(srv, ss)
	if err != nil {
		log.Printf(" [ERROR] RPC failed with error %v\n", err)
	}
	return err
}
