package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound         = sqlx.ErrNotFound    // 未找到
	ErrUserAlreadyExist = errors.New("用户已存在") // 用户已存在
)
