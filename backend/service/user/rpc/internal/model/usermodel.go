package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindOneByEmailOrMobile(context context.Context, email, mobile string) (*User, error)
		Users(ctx context.Context, cursor int64, lastLogin string, sortField string, limit int64) ([]*User, error)
		DeleteUserById(ctx context.Context, id int64) error
		UpdateLastLoginTime(ctx context.Context, data *User) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *defaultUserModel) FindOneByEmailOrMobile(ctx context.Context, email, mobile string) (*User, error) {
	var resp User
	query := "SELECT " + userRows + " FROM " + m.table + " WHERE email = ? OR mobile = ? LIMIT 1"
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, email, mobile)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Users 查询用户列表，根据指定的排序字段和限制数量返回用户信息
// 该方法支持分页和排序，可以根据用户ID或最后登录时间降序排列
// 参数:
//
//	ctx - 上下文，用于传递请求范围的上下文信息
//	cursor - 游标，用于分页查询的起始位置
//	lastLogin - 最后登录时间，当排序字段为"last_login"时使用
//	sortField - 排序字段，支持"id"和"last_login"
//	limit - 每页限制返回的用户数量
//
// 返回值:
//
//	[]*User - 用户列表
//	error - 错误信息，如果查询失败则返回错误
func (m *defaultUserModel) Users(ctx context.Context, cursor int64, lastLogin string, sortField string, limit int64) ([]*User, error) {

	var (
		resp     []*User
		query    string
		anyField any
	)

	// 根据排序字段构建查询语句
	switch sortField {
	case "id":
		query = "SELECT " + userRows + " FROM " + m.table + " WHERE id < ? AND is_deleted = 0 ORDER BY id DESC LIMIT ?" // 按照用户ID降序
		anyField = cursor
	case "last_login":
		query = "SELECT " + userRows + " FROM " + m.table + " WHERE last_login < ? AND is_deleted = 0 ORDER BY last_login DESC LIMIT ?" // 按照最后登录时间降序
		anyField = lastLogin
	default:
		query = "SELECT " + userRows + " FROM " + m.table + " WHERE id < ? AND is_deleted = 0 ORDER BY id DESC LIMIT ?" // 默认按照用户ID降序
	}

	// 执行查询，使用参数化查询来避免 SQL 注入
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, anyField, limit)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound // 无数据时返回特定错误
		}
		return nil, err // 返回其他错误
	}

	return resp, nil
}

func (m *defaultUserModel) DeleteUserById(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	dataverseUserUserEmailKey := fmt.Sprintf("%s%v", cacheDataverseUserUserEmailPrefix, data.Email)
	dataverseUserUserIdKey := fmt.Sprintf("%s%v", cacheDataverseUserUserIdPrefix, data.Id)
	dataverseUserUserMobileKey := fmt.Sprintf("%s%v", cacheDataverseUserUserMobilePrefix, data.Mobile)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, "is_deleted = 1")
		return conn.ExecCtx(ctx, query, id)
	}, dataverseUserUserEmailKey, dataverseUserUserIdKey, dataverseUserUserMobileKey)
	return err
}

func (m *defaultUserModel) UpdateLastLoginTime(ctx context.Context, data *User) error {
	dataverseUserUserEmailKey := fmt.Sprintf("%s%v", cacheDataverseUserUserEmailPrefix, data.Email)
	dataverseUserUserIdKey := fmt.Sprintf("%s%v", cacheDataverseUserUserIdPrefix, data.Id)
	dataverseUserUserMobileKey := fmt.Sprintf("%s%v", cacheDataverseUserUserMobilePrefix, data.Mobile)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, "last_login = ?")
		return conn.Exec(query, data.LastLogin, data.Id)
	}, dataverseUserUserEmailKey, dataverseUserUserIdKey, dataverseUserUserMobileKey)
	return err
}
