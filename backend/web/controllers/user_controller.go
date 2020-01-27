package controllers

import (
	"Kapi/config"
	"Kapi/services"
	"Kapi/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"net/http"
	"time"
)

type UserController struct {
	Ctx iris.Context
	UserService services.IUserService
}

func (uc *UserController) PostLogin() mvc.View {
	email := uc.Ctx.FormValue("email")
	password := uc.Ctx.FormValue("password")
	token, err := uc.UserService.LoginByEmail(email,password)
	if err != nil {
		uc.Ctx.Values().Set("err_msg",err)
		uc.Ctx.Application().Logger().Errorf("Error in PostEmail:%v",err)
		uc.Ctx.StatusCode(401)
	}

	AppConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
	if err != nil {
		uc.Ctx.Values().Set("err_msg",err)
		uc.Ctx.Application().Logger().Errorf("Error in PostEmail:%v",err)
		uc.Ctx.StatusCode(401)
		//return mvc.View{
		//	Name:"public/error.html",
		//	Data:iris.Map{
		//		"Message": fmt.Sprint(err),
		//	},
		//}
	}

	cookie := &http.Cookie{
		Name:"token",
		Value:token,
		Expires:time.Now().Add(time.Duration(AppConfig.SecurityExpiresIn)),
	}
	uc.Ctx.SetCookie(cookie)

	return mvc.View{
		Name:"public/err.html",
		Data:iris.Map{
			"Message": token,
		},
	}
}