package encrypt

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCodec 是一个模拟的 codec 包。
type MockCodec struct {
	mock.Mock
}

// EcbEncrypt 模拟加密过程。
func (m *MockCodec) EcbEncrypt(key, data []byte) ([]byte, error) {
	args := m.Called(key, data)
	return args.Get(0).([]byte), args.Error(1)
}

// EcbDecrypt 模拟解密过程。
func (m *MockCodec) EcbDecrypt(key, data []byte) ([]byte, error) {
	args := m.Called(key, data)
	return args.Get(0).([]byte), args.Error(1)
}

// 使用模拟的 codec 包。
var mcodec *MockCodec
var (
	aesKey = "DQC0DHTABQ6VQCXJ7HM6P38SFD46LPOA"
	salt   = "COGhXZR45qVOmFBz0hWza8cJw3RhX0My"
)

func TestEncMobile_Success(t *testing.T) {
	// 设置
	mcodec = new(MockCodec)
	mobile := "18623334444"
	encryptedMobile := "GR+15Nn8NUgoDJE25lvKAA=="
	mcodec.On("EcbEncrypt", []byte(aesKey), []byte(mobile)).Return([]byte(encryptedMobile), nil)
	encrypted, err := EncMobile(mobile, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, encryptedMobile, encrypted)
}

func TestEncMobile_EncryptionError(t *testing.T) {
	// 设置
	mcodec = new(MockCodec)
	mobile := ""
	mcodec.On("EcbEncrypt", []byte(aesKey), []byte(mobile)).Return(nil, errors.New("encryption error"))
	encrypted, err := EncMobile(mobile, aesKey)
	assert.Error(t, err)
	assert.Equal(t, "", encrypted)
}

func TestDecMobile_Success(t *testing.T) {
	// 设置
	mcodec = new(MockCodec)
	mobile := "18623334444"
	encryptedMobile := "GR+15Nn8NUgoDJE25lvKAA=="
	originalData := []byte(mobile)
	decodedData, _ := base64.StdEncoding.DecodeString(encryptedMobile)

	mcodec.On("EcbDecrypt", []byte(aesKey), decodedData).Return(originalData, nil)
	decrypted, err := DecMobile(encryptedMobile, aesKey)
	assert.NoError(t, err)
	assert.Equal(t, mobile, decrypted)
}

func TestDecMobile_DecryptionError(t *testing.T) {
	// 设置
	mcodec = new(MockCodec)
	encryptedMobile := ""
	decodedData, _ := base64.StdEncoding.DecodeString(encryptedMobile)

	mcodec.On("EcbDecrypt", []byte(aesKey), decodedData).Return(nil, errors.New("decryption error"))
	decrypted, err := DecMobile(encryptedMobile, aesKey)
	assert.Error(t, err)
	assert.Equal(t, "", decrypted)
}

func TestDecMobile_Base64DecodeError(t *testing.T) {
	// 设置
	invalidBase64Data := "invalidBase64Data"

	decrypted, err := DecMobile(invalidBase64Data, aesKey)
	assert.Error(t, err)
	assert.Equal(t, "", decrypted)
}

func TestVerifyPassword(t *testing.T) {
	// 创建一个测试哈希密码
	password := "GyKUWbXeq9LdYLE"
	salt := "COGhXZR45qVOmFBz0hWza8cJw3RhX0My"
	hashedPassword, err := EncPassword(password, salt)
	t.Logf("Hashed password: %s", hashedPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 测试用例：密码匹配
	if !VerifyPassword(string(hashedPassword), password, salt) {
		t.Error("VerifyPassword should return true for matching password")
	}

	// 测试用例：密码不匹配
	if VerifyPassword(string(hashedPassword), "wrongPassword", salt) {
		t.Error("VerifyPassword should return false for non-matching password")
	}
}
