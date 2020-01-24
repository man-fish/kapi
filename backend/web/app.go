package main

import (
	"Kapi/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func main() {
	f, err := utils.NewLogFile()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	/* 初始化日志文件 */
	app := iris.New()
	/* 主角登场 */
	customLogger := logger.New(logger.Config{
		//状态显示状态代码
		Status: true,
		// IP显示请求的远程地址
		IP: true,
		//方法显示http方法
		Method: true,
		// Path显示请求路径
		Path: true,
		// Query将url查询附加到Path。
		Query: true,
		Columns:true,
		// 如果不为空然后它的内容来自`ctx.Values(),Get("logger_message")
		//将添加到日志中。
		MessageContextKeys: []string{"logger_message"},
		//如果不为空然后它的内容来自`ctx.GetHeader（“User-Agent”）
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(customLogger)
	app.Logger().SetLevel("warning")
	app.Logger().SetOutput(f)
	/* 日志初始化 */
	app.StaticWeb("/static","./assets")
	/* 静态文件初始化 */
	utils.InitErrorHandler(app)
	/* 全局错误处理初始化 */
	conn, err := utils.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	/* 数据库连接初始化 */

	err = app.Run(iris.Addr(":8080"),
		iris.WithPathEscape,
		iris.WithOptimizations,
		iris.WithCharset("utf-8"),
		iris.WithTimeFormat("2006-01-02 15:04:05"),
		iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
	/* 启动服务 */
}
