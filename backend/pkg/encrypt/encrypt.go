package encrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/codec"
	"golang.org/x/crypto/bcrypt"
)

func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	return in[:]
}

// EncPassword 使用 bcrypt 对密码进行加密
func EncPassword(password, slat string) (string, error) {
	// 去除密码首尾的空白字符并拼接加密种子
	trimmedPassword := strings.TrimSpace(password + slat)
	// 使用 bcrypt 对密码进行加密，默认成本为 14
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(trimmedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(encPassword, password, slat string) bool {
	// 去除密码首尾的空白字符并拼接加密种子
	trimmedPassword := strings.TrimSpace(password + slat)
	err := bcrypt.CompareHashAndPassword([]byte(encPassword), []byte(trimmedPassword))
	return err == nil
}

func EncMobile(mobile, aesKey string) (string, error) {
	if len(mobile) == 0 {
		return "", errors.New("mobile is empty")
	}
	data, err := codec.EcbEncrypt([]byte(aesKey), []byte(mobile))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecMobile(mobile, aesKey string) (string, error) {
	if len(mobile) == 0 {
		return "", errors.New("mobile is empty")
	}
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(aesKey), originalData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
