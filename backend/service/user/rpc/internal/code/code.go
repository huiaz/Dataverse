package code

import "dataverse/pkg/xcode"

var (
	ErrUserNameEmpty    = xcode.New(100001, "用户名不能为空")
	ErrUserIdEmpty      = xcode.New(100002, "用户ID不能为空")
	ErrPasswordIsEmpty  = xcode.New(100003, "密码不能为空")
	ErrMobileIsValid    = xcode.New(100004, "手机号不合法")
	ErrEmailIsValid     = xcode.New(100005, "邮箱不合法")
	ErrUserAlreadyExist = xcode.New(100006, "用户已存在")
	ErrInvalidRequest   = xcode.New(100007, "无效的请求, 请检查请求参数")
	ErrUerIdValid       = xcode.New(100008, "用户ID参数不合法")
	ErrUserNotExist     = xcode.New(100009, "用户不存在")
	ErrPasswordIsWrong  = xcode.New(100010, "密码错误")
	ErrInternalError    = xcode.New(100011, "内部错误")
	ErrInvalidSortType  = xcode.New(100012, "无效的排序类型")
)
