package grpc_proxy_middleware

import (
	"fmt"
	"go_gateway/common"
	"go_gateway/dao"
	"log"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// GRPCFlowLimitMiddleware grpc 限流
func GRPCFlowLimitMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 服务限流
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := common.FlowLimiterHandler.GetLimiter(
				common.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}
		// 客户端限流
		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			// 取出IP
			peerCtx, ok := peer.FromContext(ss.Context())
			if !ok {
				return errors.New("peer not found with context")
			}
			peerAddr := peerCtx.Addr.String()
			addrPos := strings.LastIndex(peerAddr, ":")
			clientIP := peerAddr[0:addrPos]
			// 检测
			clientLimiter, err := common.FlowLimiterHandler.GetLimiter(
				common.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("%v flow limit %v \n", clientIP, serviceDetail.AccessControl.ClientIPFlowLimit))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("[ERROR] GrpcFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
