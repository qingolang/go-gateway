package grpc_proxy_middleware

import (
	"go_gateway/common"
	"go_gateway/dao"
	"log"

	"google.golang.org/grpc"
)

// GRPCFlowCountMiddleware 流量统计 返回ERROR则中断
func GRPCFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 增加平台
		totalCounter, err := common.FlowCounterHandler.GetCounter(common.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()
		// 增加当前服务
		serviceCounter, err := common.FlowCounterHandler.GetCounter(common.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := handler(srv, ss); err != nil {
			log.Printf("[ERROR] GrpcFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
