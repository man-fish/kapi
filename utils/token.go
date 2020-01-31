package utils

import (
	"Kapi/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyCustomClaims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

func DefaultToken(uid int64) (tokenString string, err error) {
	AppConfig, err := config.GetConfig(RootPath()+"/config/config.json")
	claims := MyCustomClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix()+int64(AppConfig.SecurityExpiresIn),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(AppConfig.SecurityKey))
	return
}

func VerifyToken(tokenString string, sercetKey string) (uid int64, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sercetKey), nil
	})
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		uid = claims.Uid
		if claims.VerifyExpiresAt(time.Now().Unix(),true) {
			return
		}else{
			err = errors.New("token已经过期啦。")
			return 0, nil
		}
	} else {
		return
	}
}

