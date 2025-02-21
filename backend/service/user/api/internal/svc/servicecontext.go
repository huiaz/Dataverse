package svc

import (
	"dataverse/pkg/interceptors"
	"dataverse/service/user/api/internal/config"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRPC  userclient.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRPCConn := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	return &ServiceContext{
		Config:   c,
		UserRPC:  userclient.NewUser(userRPCConn),
		BizRedis: rds,
	}
}
