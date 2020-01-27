package main

import (
	"Kapi/config"
	"Kapi/repositories"
	"Kapi/services"
	"Kapi/utils"
	"Kapi/backend/web/controllers"
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/mvc"
)

func main() {
	f, err := utils.NewLogFile()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	/* 初始化日志文件 */
	AppConfig, err := config.GetConfig(utils.RootPath()+"/config/config.json")
	if err != nil {
		panic(err)
	}
	/* 读取配置文件 */
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
	app.Logger().SetLevel("info")
	//app.Logger().SetOutput(f)
	/* 日志初始化 */
	template := iris.HTML("./backend/web/views",".html")
	template.Reload(true)
	app.RegisterView(template)
	/* 模版文件初始化 */
	app.StaticWeb("/static",AppConfig.StaticPath)
	/* 静态文件初始化 */
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message",ctx.Values().Get("err_msg"))
		//ctx.ViewLayout("")
		err := ctx.View("public/err.html")
		if err != nil {
			panic(err)
		}
	})
	/* 全局错误处理初始化 */
	conn, err := utils.NewMysqlConn(AppConfig.MysqlDsn)
	if err != nil {
		panic(err)
	}
	/* 数据库连接初始化 */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	/* 获取ctx结构体 */
	userRespository := repositories.NewUserManager(conn,"user")
	userService := services.NewUserService(userRespository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	/* user路由注册 */
	err = app.Run(iris.Addr(AppConfig.Port),
		iris.WithPathEscape,
		iris.WithOptimizations,
		iris.WithCharset("utf-8"),
		iris.WithTimeFormat("2006-01-02 15:04:05"),
		//iris.WithoutServerError(iris.ErrServerClosed)
	)
	if err != nil {
		panic(err)
	}
	/* 启动服务 */
}
