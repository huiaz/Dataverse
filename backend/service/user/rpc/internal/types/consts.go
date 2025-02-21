package types

const (
	DefaultPageSize       = 10
	MaxPageSize           = 500
	DefaultSortUserCursor = 1 << 30 // 默认排序用户游标
)

const (
	SortTypeUserId = iota
	SortTypeLastLogin
)
