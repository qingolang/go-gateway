package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GRPCJWTFlowCountMiddleware 租户流量统计与限流
func GRPCJWTFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		appInfos := md.Get("app")
		if len(appInfos) == 0 {
			if err := handler(srv, ss); err != nil {
				log.Printf(" [ERROR] RPC failed with error %v\n", err)
				return err
			}
			return nil
		}

		appInfo := &dao.APP{}
		if err := json.Unmarshal([]byte(appInfos[0]), appInfo); err != nil {
			return err
		}
		appCounter, err := common.FlowCounterHandler.GetCounter(common.FlowAppPrefix + appInfo.APPID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			return errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.QPD, appCounter.TotalCount))
		}
		if err := handler(srv, ss); err != nil {
			log.Printf(" [ERROR] RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
