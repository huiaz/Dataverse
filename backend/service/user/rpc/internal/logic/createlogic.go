package logic

import (
	"context"
	"regexp"
	"time"

	"dataverse/pkg/encrypt"
	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/model"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func vilidateEmail(email string) bool {
	// 定义邮箱的正则表达式
	// 1. 邮箱名称：由字母、数字、点、下划线和中划线组成，长度为1-64个字符
	// 2. @符号
	// 3. 域名：由字母、数字、点和中划线组成，长度为1-255个字符
	if len(email) == 0 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)

}
func validateMobile(mobile string) bool {
	// 定义手机号的正则表达式
	// 1. 中国的手机号：11位数字，以1开头
	// 2. 国际手机号：以+号开头，后跟国家代码和手机号
	if len(mobile) == 0 {
		return false
	}
	mobileRegex := regexp.MustCompile(`^(?:\+?\d{1,3}[-.\s]?)?(?:\d{1,4}[-.\s]?)?\d{7,14}$`)
	return mobileRegex.MatchString(mobile)
}

func (l *CreateLogic) Create(in *user.CreateRequest) (*user.CreateResponse, error) {
	// 参数校验
	if len(in.Name) == 0 {
		return nil, code.ErrUserNameEmpty
	}
	if len(in.Password) == 0 {
		return nil, code.ErrPasswordIsEmpty
	}
	if !validateMobile(in.Mobile) {
		return nil, code.ErrMobileIsValid
	}
	if !vilidateEmail(in.Email) {
		return nil, code.ErrEmailIsValid
	}
	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}

	// 创建用户
	password, err := encrypt.EncPassword(in.Password, l.svcCtx.Config.PasswordSalt)
	if err != nil {
		logx.Errorf("密码加密异常:EncPassword: error: %v", err)
		return nil, code.ErrInternalError
	}
	newUser := &model.User{
		Username:   in.Name,
		Password:   password,
		Email:      in.Email,
		Mobile:     in.Mobile,
		IsAdmin:    in.IsAdmin,
		IsDeleted:  false,
		CreateTime: time.Now().Local(),
	}

	// 查询用户是否存在
	userExist, err := l.svcCtx.UserModel.FindOneByEmailOrMobile(l.ctx, in.Email, in.Mobile)

	switch err {
	case nil:
		if userExist.IsDeleted {
			// 更新用户状态
			userExist.IsDeleted = false
			userExist.Username = in.Name
			userExist.Password = password
			userExist.IsAdmin = in.IsAdmin
			userExist.Mobile = in.Mobile
			userExist.Email = in.Email
			userExist.UpdateTime = time.Now().Local()
			err := l.svcCtx.UserModel.Update(l.ctx, userExist)
			if err != nil {
				logx.Errorf("对已删除的用户进行恢复,更新操作失败:Update: req: %v, error: %v", in, err)
				return nil, code.ErrInternalError
			}
			return &user.CreateResponse{UserId: userExist.Id}, nil
		}
		return nil, code.ErrUserAlreadyExist
	case model.ErrNotFound:
		ret, err := l.svcCtx.UserModel.Insert(l.ctx, newUser)
		if err != nil {
			logx.Errorf("数据库新建用户异常:Insert: req: %v, error: %v", in, err)
			return nil, code.ErrInternalError
		}
		userId, err := ret.LastInsertId()
		if err != nil {
			logx.Errorf("获取数据库新建用户ID异常:LastInsertId: error: %v", err)
			return nil, code.ErrInternalError
		}
		return &user.CreateResponse{UserId: userId}, nil
	default:
		logx.Errorf("数据库查询用户异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}
}
