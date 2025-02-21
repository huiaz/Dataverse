package logic

import (
	"context"

	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/model"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	// 输入验证
	if in == nil || in.UserId <= 0 || in.UserId != int64(int(in.UserId)) {
		return nil, code.ErrUerIdValid
	}

	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}

	resp, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, code.ErrUserNotExist
		}
		logx.Errorf("数据库查询用户异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}

	// 确保 resp 不为空
	if resp == nil {
		logx.Errorf("数据库查询用户返回空结果: req: %v", in)
		return nil, code.ErrInternalError
	}

	return &user.UserInfoResponse{
		UserId:     resp.Id,
		Name:       resp.Username,
		Email:      resp.Email,
		Mobile:     resp.Mobile,
		IsAdmin:    resp.IsAdmin,
		IsDelete:   resp.IsDeleted,
		LastLogin:  resp.LastLogin.Time.Unix(),
		CreateTime: resp.CreateTime.Unix(),
		UpdateTime: resp.UpdateTime.Unix(),
	}, nil
}
