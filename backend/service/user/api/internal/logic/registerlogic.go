package logic

import (
	"context"
	"dataverse/service/user/api/internal/code"
	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	if !vilidateEmail(req.Email) {
		return nil, code.ErrEmailIsValid
	}
	if len(req.Password) == 0 {
		return nil, code.ErrPasswordIsEmpty
	}
	if !validateMobile(req.Mobile) {
		return nil, code.ErrMobileIsValid
	}
	if len(req.Name) == 0 {
		return nil, code.ErrUserNameEmpty
	}

	// 调用 RPC 的 Register 方法进行注册
	userInfo, err := l.svcCtx.UserRPC.Create(l.ctx, &userclient.CreateRequest{
		IsAdmin:  req.IsAdmin,
		Email:    req.Email,
		Password: req.Password,
		Mobile:   req.Mobile,
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}
	// 构建响应
	resp = &types.RegisterResponse{
		UserId: userInfo.UserId,
	}

	return
}
