package logic

import (
	"context"

	"dataverse/service/user/api/internal/code"
	"dataverse/service/user/api/internal/svc"
	"dataverse/service/user/api/internal/types"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteRequest) (resp *types.DeleteResponse, err error) {
	if req.UserId <= 0 {
		return nil, code.ErrUerIdValid
	}
	data, err := l.svcCtx.UserRPC.Delete(l.ctx, &user.DeleteRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.DeleteResponse{
		UserId: data.UserId,
	}

	return resp, nil
}
