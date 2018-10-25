package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	Secret_Key = "this is wechatpay !"
)

func TokenGenerate(userId int64, expireTime int64) (token string, err error) {

	ojwt := jwt.New(jwt.SigningMethodHS256)

	ojwt.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(expireTime)).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    fmt.Sprintf("%v", userId),
	}

	return ojwt.SignedString([]byte(Secret_Key))
}

func TokenValidate(token string) (userId int64, err error) {

	ojwt, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret_Key), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := ojwt.Claims.(*jwt.StandardClaims); ok && ojwt.Valid {

		userId, _ = strconv.ParseInt(claims.Issuer, 10, 64)

		return userId, nil
	}

	return 0, errors.New("token valide failse,please login again")

}
