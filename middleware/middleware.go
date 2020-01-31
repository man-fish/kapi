package middleware

import (
	"Kapi/config"
	"Kapi/utils"
	"github.com/kataras/iris"
)

func AuthWithToken(ctx iris.Context) {
	signature := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjE0LCJleHAiOjE1ODA5MjAwOTF9.UzLS63eauY7UNCMcPLcJ0vOCG_AJsB3Ls3tGRjZsSDw"//ctx.GetCookie("token")
	appConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
	if err != nil {
		utils.ErrorWithCode(err,"auth",500,ctx)
		return
	}
	uid, err := utils.VerifyToken(signature,appConfig.SecurityKey)
	if err != nil {
		utils.ErrorWithCode(err,"auth",401,ctx)
		return
	}
	ctx.Values().Set("uid",uid)
	ctx.Next()
}