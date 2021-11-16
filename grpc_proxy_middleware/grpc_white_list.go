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

// GRPCWhiteListMiddleware 匹配接入方式 基于请求信息
func GRPCWhiteListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		if serviceDetail.AccessControl.OpenWhiteList != 1 {
			if err := handler(srv, ss); err != nil {
				log.Printf(" [ERROR] RPC failed with error %v\n", err)
				return err
			}
			return nil
		}

		// 取出IP
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		// 校验白名单
		if serviceDetail.AccessControl.WhiteList != "" {
			if !common.InStringSlice(strings.Split(serviceDetail.AccessControl.WhiteList, ","), clientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf(" [ERROR] RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
