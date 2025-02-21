package svc

import (
	"dataverse/service/user/rpc/internal/config"
	"dataverse/service/user/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config             config.Config
	UserModel          model.UserModel
	GroupModel         model.GroupModel
	UserGroupUserModel model.UserGroupModel
	UserRedis          *redis.Redis
	SingleFlightGroup  singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	return &ServiceContext{
		Config:             c,
		UserModel:          model.NewUserModel(conn, c.CacheRedis),
		GroupModel:         model.NewGroupModel(conn, c.CacheRedis),
		UserGroupUserModel: model.NewUserGroupModel(conn, c.CacheRedis),
		UserRedis:          rds,
	}
}
