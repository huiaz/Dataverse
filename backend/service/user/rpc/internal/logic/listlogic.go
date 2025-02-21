package logic

import (
	"context"
	"time"

	"dataverse/service/user/rpc/internal/code"
	"dataverse/service/user/rpc/internal/svc"
	"dataverse/service/user/rpc/internal/types"
	"dataverse/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// validateSortType 验证排序类型是否有效
func validateSortType(sortType int32) bool {
	return sortType == types.SortTypeUserId || sortType == types.SortTypeLastLogin
}

// setDefaultPagination 设置默认分页参数
func setDefaultPagination(in *user.ListRequest) {
	// 设置默认分页参数
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize // 默认页面大小
	} else if in.PageSize > types.MaxPageSize {
		in.PageSize = types.MaxPageSize // 限制最大页面大小
	}
	if in.Cursor == 0 {
		switch in.SortType {
		case types.SortTypeLastLogin:
			in.Cursor = time.Now().Unix()
		default:
			in.Cursor = types.DefaultSortUserCursor
		}
	}
}

// List 实现了用户列表的查询逻辑。
// 它接受一个 ListRequest 作为输入，返回一个 ListResponse 和一个错误对象。
// 该函数主要负责根据输入参数中的排序类型和分页信息，从数据库中查询用户列表，并构造响应对象。
func (l *ListLogic) List(in *user.ListRequest) (*user.ListResponse, error) {
	// 参数校验
	if !validateSortType(in.SortType) {
		logx.Errorf("无效的排序类型: req: %v", in)
		return nil, code.ErrInvalidSortType
	}

	// 设置默认分页参数
	setDefaultPagination(in)
	var (
		sortField string
		err       error
		sortId    int64
		lastLogin string
	)
	// 获取排序字段
	switch in.SortType {
	case types.SortTypeLastLogin:
		sortField = "last_login"
		lastLogin = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	case types.SortTypeUserId:
		sortField = "id"
		sortId = in.Cursor
	default:
		sortField = "id"
		sortId = in.Cursor
	}

	// 检查 svcCtx 和 ctx 是否为空
	if l.svcCtx == nil || l.ctx == nil {
		logx.Errorf("svcCtx 或 ctx 为空: req: %v", in)
		return nil, code.ErrInternalError
	}

	// 查询用户列表
	resp, err := l.svcCtx.UserModel.Users(l.ctx, sortId, lastLogin, sortField, in.PageSize)
	if err != nil {
		logx.Errorf("数据库查询用户列表异常: req: %v, error: %v", in, err)
		return nil, code.ErrInternalError
	}

	var (
		curPage []*user.UserInfoResponse
		isEnd   bool = len(resp) < int(in.PageSize)
		cursor  int64
		lastId  int64
	)

	// 构造返回的用户信息
	for _, u := range resp {
		curPage = append(curPage, &user.UserInfoResponse{
			UserId:     u.Id,
			Name:       u.Username,
			Email:      u.Email,
			Mobile:     u.Mobile,
			IsAdmin:    u.IsAdmin,
			IsDelete:   u.IsDeleted,
			LastLogin:  u.LastLogin.Time.Unix(),
			CreateTime: u.CreateTime.UTC().Unix(),
			UpdateTime: u.UpdateTime.UTC().Unix(),
		})
	}

	// 计算最后的 cursor 和 lastId
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.UserId
		if in.SortType == types.SortTypeLastLogin {
			cursor = pageLast.LastLogin
		} else {
			cursor = lastId
		}
		if cursor <= 0 {
			cursor = 0
		}
	}

	return &user.ListResponse{
		Users:  curPage,
		IsEnd:  isEnd,
		Cursor: cursor,
		UserId: lastId,
	}, nil
}
