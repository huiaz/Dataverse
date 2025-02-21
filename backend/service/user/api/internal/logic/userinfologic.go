package logic

import (
	"context"

	"dataverse/service/user/api/internal/code"
	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// 校验参数
	if req == nil || req.UserId <= 0 || req.UserId != int64(int(req.UserId)) {
		return nil, code.ErrUerIdValid
	}
	// 调用 RPC 的 UserInfo 方法获取用户信息
	u, err := l.svcCtx.UserRPC.UserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UserInfoResponse{
		UserId:     u.UserId,
		Name:       u.Name,
		Email:      u.Email,
		Mobile:     u.Mobile,
		IsAdmin:    u.IsAdmin,
		IsDelete:   u.IsDelete,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}

	return
}
