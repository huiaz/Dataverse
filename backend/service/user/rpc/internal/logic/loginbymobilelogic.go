package logic

import (
	"context"
	"database/sql"
	"time"

	"dataverse/pkg/encrypt"
	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/model"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileLogic {
	return &LoginByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginByMobileLogic) LoginByMobile(in *user.LoginByMobileRequest) (*user.LoginResponse, error) {
	if !validateMobile(in.Mobile) {
		return nil, code.ErrMobileIsValid
	}
	if len(in.Password) == 0 {
		return nil, code.ErrPasswordIsEmpty
	}
	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}
	resp, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.ErrUserNotExist
		}
		logx.Errorf("数据库查询用户异常:Mobile: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}
	// 校验密码
	if !encrypt.VerifyPassword(resp.Password, in.Password, l.svcCtx.Config.PasswordSalt) {
		return nil, code.ErrPasswordIsWrong
	}
	// 使用异步方式更新用户最后登录时间
	go l.updateLastLoginTime(resp)

	return &user.LoginResponse{UserId: resp.Id}, nil
}

func (l *LoginByMobileLogic) updateLastLoginTime(user *model.User) {
	// 修复 context.WithTimeout 的用法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user.LastLogin = sql.NullTime{Time: time.Now(), Valid: true}
	if err := l.svcCtx.UserModel.UpdateLastLoginTime(ctx, user); err != nil {
		if ctx.Err() == context.Canceled {
			logx.Errorf("数据库更新用户最后登录时间异常:Update:, error: context canceled, reason: %v", ctx.Err())
		} else {
			logx.Errorf("数据库更新用户最后登录时间异常:Update:, error: %v", err)
		}
	}
}
