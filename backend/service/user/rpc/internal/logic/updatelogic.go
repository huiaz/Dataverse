package logic

import (
	"context"

	"dataverse/pkg/encrypt"
	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/model"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLogic) Update(in *user.UpdateRequest) (*user.UpdateResponse, error) {
	// 输入验证
	if in == nil || in.UserId <= 0 || in.UserId != int64(int(in.UserId)) {
		return nil, code.ErrUerIdValid
	}
	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}

	// 尝试获取用户
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.ErrUserNotExist
		}
		logx.Errorf("数据库查询用户异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}
	if in.Password != "" {
		// 密码加密
		passwd, err := encrypt.EncPassword(in.Password, l.svcCtx.Config.PasswordSalt)
		if err != nil {
			logx.Errorf("密码加密异常:EncPassword: error: %v", err)
			return nil, code.ErrInternalError
		}
		userInfo.Password = passwd
	}
	if in.Email != "" {
		userInfo.Email = in.Email
	}
	if in.Mobile != "" {
		userInfo.Mobile = in.Mobile
	}
	if in.Name != "" {
		userInfo.Username = in.Name
	}
	if in.IsAdmin != userInfo.IsAdmin {
		userInfo.IsAdmin = in.IsAdmin
	}
	// 更新用户
	if err := l.svcCtx.UserModel.Update(l.ctx, userInfo); err != nil {
		logx.Errorf("数据库更新用户异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}
	return &user.UpdateResponse{UserId: in.UserId}, nil

}
