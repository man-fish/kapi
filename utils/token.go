package utils

import (
	"Kapi/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func DefaultToken(uid int64) (tokenString string, err error) {
	AppConfig, err := config.GetConfig(RootPath()+"/config/config.json")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"nbf": time.Now().Unix() + int64(AppConfig.SecurityExpiresIn),
	})
	tokenString, err = token.SignedString(AppConfig.SecurityKey)
	return
}
