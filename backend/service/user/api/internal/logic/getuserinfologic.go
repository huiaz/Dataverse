package logic

import (
	"context"
	"encoding/json"

	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResponse, err error) {
	// 从 svcCtx 中获取用户信息
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return &types.UserInfoResponse{}, nil
	}
	userInfo, err := l.svcCtx.UserRPC.UserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UserInfoResponse{
		UserId:     userInfo.UserId,
		Name:       userInfo.Name,
		Mobile:     userInfo.Mobile,
		Email:      userInfo.Email,
		IsAdmin:    userInfo.IsAdmin,
		CreateTime: userInfo.CreateTime,
		UpdateTime: userInfo.UpdateTime,
		LastLogin:  userInfo.LastLogin,
		IsDelete:   userInfo.IsDelete,
	}
	return
}
