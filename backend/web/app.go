package main

import (
	"Kapi/backend/web/controllers"
	"Kapi/config"
	"Kapi/middleware"
	"Kapi/repositories"
	"Kapi/services"
	"Kapi/utils"
	"context"
	"database/sql"
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
		ctx.JSON(utils.MakeDefaultRes(0, ctx.Values().Get("err_msg") , nil))
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
	mvc.Configure(app.Party("/user"), userHandler(conn, ctx))
	mvc.Configure(app.Party("/group"), groupHandler(conn, ctx))
	/* 路由注册 */

	err = app.Run(iris.Addr(AppConfig.Port),
		iris.WithPathEscape,
		iris.WithOptimizations,
		iris.WithCharset("utf-8"),
		iris.WithTimeFormat("2006-01-02 15:04:05"))
		//iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
	/* 启动服务 */
}

func userHandler(conn *sql.DB, ctx context.Context) func(*mvc.Application) {
	return func(app *mvc.Application) {
		userRepository := repositories.NewUserManager(conn,"Kapi_user")
		userService := services.NewUserService(userRepository)
		app.Register(ctx, userService)
		app.Handle(new(controllers.UserController))
	}
}

func groupHandler(conn *sql.DB, ctx context.Context) func(*mvc.Application) {
	return func(app *mvc.Application) {
		app.Router.Use(middleware.AuthWithToken)
		userRepository := repositories.NewUserManager(conn,"Kapi_user")
		groupRepository := repositories.NewGroupManager(conn,"Kapi_group")
		memberRepository := repositories.NewMemberManager(conn, "Kapi_member")
		groupService := services.NewGroupService(userRepository,groupRepository,memberRepository)
		app.Register(ctx, groupService)
		app.Handle(new(controllers.GroupController))
	}
}
