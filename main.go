package main

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

	"github.com/biluohc/kloud/api"
	"github.com/biluohc/kloud/model"
)

func main() {
	err := model.InitDB(dbConfig)
	defer model.CloseDB()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fun()
}

func fun() {
	e := echo.New()
	e.Use(api.KloudContext)

	e.HTTPErrorHandler = api.HTTPErrorHandler
	e.Use(middleware.Recover())

	e.Use(logger)
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "[${time_rfc3339} ${remote_ip}] ${status} ${method} ${uri}\n",
	// }))

	// 本地调试前端需要跨域
	// e.Use(middleware.CORSWithConfig(
	// 	middleware.CORSConfig{
	// 		Skipper:          middleware.DefaultSkipper,
	// 		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8000"},
	// 		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	// 		AllowCredentials: true,
	// 	}))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// spa
	e.File("/", "Afront/dist/index.html")
	// statics
	e.File("/favicon.ico", "Afront/dist/static/favicon.ico")
	e.Static("/static", "Afront/dist/static/")

	// /apiv1/*
	g := e.Group("/apiv1")
	g.POST("/login", api.Login)
	g.POST("/signin", api.Signin)

	// 有密码临时重定向到前端，填写密码后post，通过验证后save cookie/token再重定向回去
	e.GET("/apiv1/share/info/:id", api.GetShareInfo)
	e.POST("/apiv1/share/init/:id", api.ShareInit)
	g.GET("/share/:id/:name", api.GetShareFile)
	g.GET("/public/:id/:name", api.GetPublicFile)

	g.Use(api.AuthRequired)
	// 使用中间件后才要求已经登录
	g.GET("/login", api.CurLogin)
	g.DELETE("/login", api.ExitLogin)
	g.POST("/user/avatar", api.UpdateAvatarURL)
	g.POST("/user/password", api.UpdatePassword)

	g.GET("/entry/info", api.GetEntryRootInfo)
	g.GET("/entry/info/", api.GetEntryRootInfo)
	g.GET("/entry/info/:id", api.GetEntryInfo)
	// 从 / 到当前id的 info 列表
	g.GET("/entry/path", api.GetEntryRootPath)
	g.GET("/entry/path/", api.GetEntryRootPath)
	g.GET("/entry/path/:id", api.GetEntryPath)
	// 当前条目子条目列表
	g.GET("/entry/subs", api.GetEntryRootSubs)
	g.GET("/entry/subs/", api.GetEntryRootSubs)
	g.GET("/entry/subs/:id", api.GetEntrySubs)

	// 有name浏览器才会用文件名保存
	g.GET("/entry/file/:id/:name", api.GetEntryFile)
	g.GET("/entry/trashes", api.GetEntryTrashes)

	// 用 json 的话就不能从 URL 看到名字，相当于一点点保护吧
	g.POST("/entry/dir/:id", api.CreateDir)
	// 上传文件
	g.POST("/entry/file/:id", api.CreateFile)
	// 根据sha上传文件，如果关闭快速上传就不要载入这个api, 用户就不能通过sha1得到文件
	g.POST("/entry/sha1/:id", api.CreateFileBySha1)

	g.PUT("/entry/rename/:id", api.EntryRename)
	g.PUT("/entry/moveto/:id", api.EntryMoveTo)
	// 分享因为表, 不能选多个条目, 那就也不要加上批量啦
	g.PUT("/entry/trash/:id", api.EntryTrash)
	g.PUT("/entry/untrash/:id", api.EntryUnTrash)
	// 删除条目
	g.DELETE("/entry/:id", api.EntryDelete)

	// 分享该条目，目前只能是文件
	// 拦截一切referer, 防止盗图
	// 可选密码，放在json
	g.POST("/entry/share/:id", api.EntryShare)
	g.GET("/shares", api.Shares)
	// 修改密码
	g.PUT("/share/password/:id", api.SharePassword)
	g.PUT("/share/maxage/:id", api.ShareMaxAge)
	g.DELETE("/share/:id", api.ShareDelete)

	// 公开该条目，目前只能是文件
	// 可选 referer 正则(go的)
	g.POST("/entry/public/:id", api.EntryPublic)
	g.GET("/publics", api.Publics)
	g.PUT("/public/referer/:id", api.PublicReferer)
	g.PUT("/public/maxage/:id", api.PublicMaxAge)
	g.DELETE("/public/:id", api.PublicDelete)

	// Admin
	gAdmin := e.Group("/Apiv1")
	gAdmin.Use(api.AuthRequired)
	gAdmin.Use(api.AdminAuthRequired)
	gAdmin.GET("/login", api.CurLogin)
	gAdmin.File("/log.txt", "log/dist.txt")
	gAdmin.DELETE("/log.txt", api.ClearLogFile)
	gAdmin.GET("/stat", api.GetStat)
	gAdmin.GET("/users", api.GetUsers)
	gAdmin.GET("/user/:id", api.GetUser)
	gAdmin.PUT("/user/:id", api.ModifyUser)
	gAdmin.PUT("/user/disable/:id", api.UserDisable)
	gAdmin.PUT("/user/undisable/:id", api.UserUnDisable)
	gAdmin.POST("/user", api.NewUser)

	gAdmin.GET("/files", api.GetFiles)
	gAdmin.GET("/file/:id/:name", api.GetFile)
	gAdmin.PUT("/file/disable/:id", api.FileDisable)
	gAdmin.PUT("/file/undisable/:id", api.FileUnDisable)

	if https {
		fmt.Println(e.StartTLS(ServeAddr, httpsCrt, httpsKey))
	} else {
		fmt.Println(e.Start(ServeAddr))
	}
}
