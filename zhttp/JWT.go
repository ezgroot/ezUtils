package zhttp

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

func secret(key string) jwt.Keyfunc {
	ValidationKeyGetter := func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}

	return ValidationKeyGetter
}

// CreateToken create token
/*
	iss: 签发者
	sub: 面向的用户
	aud: 接收方
	exp: 过期时间
	nbf: 生效时间
	iat: 签发时间
	jti: 唯一身份标识
*/
func CreateToken(payload map[string]interface{}, key string) (string, error) {
	var claim jwt.MapClaims
	claim = payload

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(key))
}

// ParseToken analyse token
func ParseToken(tokenStr string, key string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, secret(key))
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to MapClaims")
		return claim, err
	}

	if !token.Valid {
		err = errors.New("token is invalid")
		return claim, nil
	}

	return claim, nil
}
