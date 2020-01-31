package utils

import (
	"fmt"
	"github.com/kataras/iris"
)

type MyError struct {
	code	int64
	msg		string
}

func NewError(code int64, msg string) *MyError {
	return &MyError{
		code: code,
		msg:  msg,
	}
}
/* 自定义错误类型，方便service层错误码跑出 */
func (e *MyError) Error() string{
	return fmt.Sprintf("err_code:%v,err_msg:%v",e.code,e.msg)
}

func (e *MyError) Code() int64 {
	return e.code
}

func (e *MyError) Msg() string {
	return e.msg
}

func ErrorWithError(err error,err_module string,ctx iris.Context) {
	ctx.Values().Set("err_msg",err.Error())
	ctx.Application().Logger().Errorf("Error in %v:%v",err_module,err)
	ctx.StatusCode(int(err.(*MyError).Code()))
}
/* 通过自定义错误触发全局错误处理 */
func ErrorWithCode(err error,err_module string,err_code int,ctx iris.Context) {
	ctx.Values().Set("err_msg",err.Error())
	ctx.Application().Logger().Errorf("Error in %v:%v",err_module,err)
	ctx.StatusCode(err_code)
}
/* 通过httpStatusCode触发全局错误处理 */

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