package grpc_proxy_middleware

import (
	"go_gateway/common"
	"go_gateway/dao"
	"log"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GRPCJWTAuthTokenMiddleware jwt auth token
func GRPCJWTAuthTokenMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 是否开启鉴权
		if serviceDetail.AccessControl.OpenAuth != 1 {
			return nil
		}
		// 取出令牌
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}
		authToken := ""
		auths := md.Get("authorization")
		if len(auths) <= 0 {
			return errors.New("miss authorization from metadata")
		}
		authToken = auths[0]
		token := strings.ReplaceAll(authToken, "Bearer ", "")
		appMatched := false

		// 鉴权
		if token != "" {
			claims, err := common.JWTDecode(token)
			if err != nil {
				return errors.WithMessage(err, "JwtDecode")
			}
			appList := dao.APPManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.APPID == claims.Issuer {
					md.Set("app", common.Obj2Json(appInfo))
					appMatched = true
					break
				}
			}
		}
		if !appMatched {
			return errors.New("not match valid app")
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("[ERROR] GrpcJwtAuthTokenMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
