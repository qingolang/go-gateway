package reverse_proxy

import (
	"context"
	"go_gateway/reverse_proxy/load_balance"
	"log"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewGrpcLoadBalanceHandler GRPC策略
func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithDefaultCallOptions(), grpc.WithInsecure())
			if err != nil {
				return nil, nil, err
			}
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
