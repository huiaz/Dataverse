package logic

import "regexp"

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

func validatePassword(password string) bool {
	return len(password) != 0
}
