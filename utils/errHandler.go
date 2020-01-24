package utils

import "github.com/kataras/iris"

func Notfound(ctx iris.Context) {
	ctx.WriteString("error/404")
}

func IntervalServerError(ctx iris.Context) {
	ctx.WriteString("err/500")
}

func Forbidden(ctx iris.Context) {
	ctx.WriteString("err/403")
}

func BadRequest(ctx iris.Context) {
	ctx.WriteString("err/400")
}

func UnAuthorized(ctx iris.Context){
	ctx.WriteString("err/401")
}

func InitErrorHandler(app *iris.Application) {
	app.OnErrorCode(iris.StatusForbidden,Forbidden)
	app.OnErrorCode(iris.StatusNotFound,Notfound)
	app.OnErrorCode(iris.StatusUnauthorized,UnAuthorized)
	app.OnErrorCode(iris.StatusBadRequest,BadRequest)
	app.OnErrorCode(iris.StatusInternalServerError,IntervalServerError)
}