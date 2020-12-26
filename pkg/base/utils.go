package base

import (
	"crypto/md5"
	"encoding/base64"
)

// PasswordAlgo 密码加密
func PasswordAlgo(password string) string {
	return base64.StdEncoding.EncodeToString(md5.New().Sum(md5.New().Sum(md5.New().Sum([]byte(password)))))
}
