package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//SignWxToken 生成token,uid用户id，expireSec过期秒数
func SignWxToken(uid int64, expireSec int) (tokenStr string, err error) {
	// 带权限创建令牌
	claims := make(jwt.MapClaims)
	claims["uid"] = uid
	claims["admin"] = false

	sec := time.Duration(expireSec)
	claims["exp"] = time.Now().Add(time.Second * sec).Unix() //自定义有效期，过期需要重新登录获取token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用自定义字符串加密 and get the complete encoded token as a string
	tokenStr, err = token.SignedString([]byte("xxxYOUR KEYxxx"))

	return tokenStr, err
}
