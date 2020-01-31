package controllers

import (
	"Kapi/config"
	"Kapi/services"
	"Kapi/utils"
	"Kapi/validator"
	"github.com/kataras/iris"
	"gopkg.in/validator.v2"
	"net/http"
	"time"
)

type UserController struct {
	Ctx iris.Context
	UserService services.IUserService
}

func (uc *UserController) PostLogin() {
	loginValidator := new(validators.LoginValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(uc.Ctx.FormValues(),loginValidator); err != nil {
		utils.ErrorWithCode(err,"PostEmail",400,uc.Ctx)
		return
	}
	if err := validator.Validate(loginValidator); err != nil {
		utils.ErrorWithCode(err,"PostEmail",400,uc.Ctx)
		return
	}
	/* 校验器，没错这些都是校验器 */
	token, err := uc.UserService.LoginByEmail(loginValidator.Email,loginValidator.Password)
	if err != nil {
		utils.ErrorWithError(err,"PostEmail,",uc.Ctx)
		return
	}
	/* 登陆业务，返回token */
	AppConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
	if err != nil {
		utils.ErrorWithCode(err,"PostEmail",400,uc.Ctx)
		return
	}
	cookie := &http.Cookie{
		Name:"token",
		Value:token,
		Expires:time.Now().Add(time.Duration(AppConfig.SecurityExpiresIn)),
	}
	uc.Ctx.SetCookie(cookie)
	/* 获取全局配置并且设置token于Cookie */
	uc.Ctx.JSON(utils.MakeDefaultRes(1,"登陆成功！",nil))
	/* 返回 */
}

func (uc *UserController) PostRegister() {
	registeValidator := new(validators.RegisterValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(uc.Ctx.FormValues(),registeValidator); err != nil {
		utils.ErrorWithCode(err,"PostRegister",400,uc.Ctx)
		return
	}
	if err := validator.Validate(registeValidator); err != nil {
		utils.ErrorWithCode(err,"PostRegister",400,uc.Ctx)
		return
	}
	/* 校验器，没错这些都是校验器 */
	ip := uc.Ctx.RemoteAddr()
	_, err := uc.UserService.RegisterByEmail(registeValidator.Username,registeValidator.Email,registeValidator.Password,ip)
	if err != nil {
		utils.ErrorWithError(err,"PostRegister",uc.Ctx)
		return
	}
	/* 注册 */
	uc.Ctx.JSON(utils.MakeDefaultRes(1,"注册成功！",nil))
	/* 返回 */
}

