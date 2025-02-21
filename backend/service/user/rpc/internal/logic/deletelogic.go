package logic

import (
	"context"

	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/model"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLogic) Delete(in *user.DeleteRequest) (*user.DeleteResponse, error) {
	// 参数校验
	if in.UserId <= 0 {
		return nil, code.ErrUerIdValid
	}
	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}
	// 直接尝试删除用户
	err := l.svcCtx.UserModel.DeleteUserById(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.ErrUserNotExist
		}
		logx.Errorf("数据库更新用户状态异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}

	return &user.DeleteResponse{UserId: in.UserId}, nil
}
