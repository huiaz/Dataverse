package logic

import (
	"context"

	"dataverse/service/user/api/internal/code"
	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.UpdateResponse, err error) {

	if req == nil || req.UserId <= 0 || req.UserId != int64(int(req.UserId)) {
		return nil, code.ErrUerIdValid
	}
	if req.Email != "" && !vilidateEmail(req.Email) {
		return nil, code.ErrEmailIsValid
	}
	if req.Mobile != "" && !validateMobile(req.Mobile) {
		return nil, code.ErrMobileIsValid
	}
	// 调用 RPC 的 Update 方法进行更新
	u, err := l.svcCtx.UserRPC.Update(l.ctx, &userclient.UpdateRequest{
		UserId:   req.UserId,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Name:     req.Name,
		Password: req.Password,
		IsAdmin:  req.IsAdmin,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.UpdateResponse{
		UserId: u.UserId,
	}
	return
}
