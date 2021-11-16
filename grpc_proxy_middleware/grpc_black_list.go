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

// GRPCBlackListMiddleware 匹配接入方式 基于请求信息
func GRPCBlackListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if serviceDetail.AccessControl.OpenBlackList != 1 {
			if err := handler(srv, ss); err != nil {
				log.Printf(" [ERROR] RPC failed with error %v\n", err)
				return err
			}
			return nil
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if serviceDetail.AccessControl.WhiteList != "" {
			// 如果存在与白名单则不校验
			for _, whileIp := range strings.Split(serviceDetail.AccessControl.WhiteList, ",") {
				if clientIP == whileIp {
					if err := handler(srv, ss); err != nil {
						log.Printf(" [ERROR] RPC failed with error %v\n", err)
						return err
					}
					return nil
				}
			}
		}

		// 校验黑名单
		if serviceDetail.AccessControl.BlackList != "" {
			if common.InStringSlice(strings.Split(serviceDetail.AccessControl.BlackList, ","), clientIP) {
				return errors.New(fmt.Sprintf("%s in black ip list", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf(" [ERROR] RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
