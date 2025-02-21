package logic

import (
	"context"

	"dataverse/pkg/jwt"
	"dataverse/pkg/xcode"
	"dataverse/service/user/api/internal/code"
	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {

	// 根据 req.Email 或者 req.Mobile 登录
	var (
		loginType string
		userInfo  *userclient.LoginResponse
	)
	// 判断登录类型
	if req.Email != "" {
		loginType = "email"
	} else if req.Mobile != "" {
		loginType = "mobile"
	} else {
		return nil, code.ErrInvalidRequest
	}
	if !validatePassword(req.Password) {
		return nil, code.ErrPasswordIsEmpty
	}

	switch loginType {
	case "email":
		// 调用 RPC 的 LoginByEmail 方法进行登录
		userInfo, err = l.svcCtx.UserRPC.LoginByEmail(l.ctx, &userclient.LoginByEmailRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return nil, err
		}
	case "mobile":
		// 调用 RPC 的 LoginByMobile 方法进行登录
		userInfo, err = l.svcCtx.UserRPC.LoginByMobile(l.ctx, &userclient.LoginByMobileRequest{
			Mobile:   req.Mobile,
			Password: req.Password,
		})
		if err != nil {
			return nil, err
		}
	}
	if userInfo == nil || userInfo.UserId == 0 {
		return nil, xcode.AccessDenied
	}
	// 构建 token
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": userInfo.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp = &types.LoginResponse{
		UserId: userInfo.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}

	return
}
