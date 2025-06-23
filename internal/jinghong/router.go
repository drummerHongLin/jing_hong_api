package jinghong

import (
	"jonghong/internal/jinghong/controller/v1/ali"
	"jonghong/internal/jinghong/controller/v1/email"
	"jonghong/internal/jinghong/controller/v1/user"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	emailservice "jonghong/internal/pkg/emailservice"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/middleware"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 采用gin框架搭建服务路由
// 路由中可以添加中间件，路由及中间件的调用顺序和添加顺序相关
// 常用的中间件有 验证身份和授权

func initRouter(g *gin.Engine) error {

	// 没有对应功能
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.InternalServerError, nil)
	})

	template := filepath.Join(known.HomeDir, "/static/html/*")

	// 静态html的地址
	g.LoadHTMLGlob(template)

	// 初始化controller
	uc := user.NewUserController(store.S)
	ac := ali.NewAliController(store.S)
	ec := email.NewEmailController(store.S, emailservice.MS)
	// 登录不需要token验证，采用用户名和密码验证
	// 回调函数也不需要
	g.POST("/login", uc.Login)

	g.GET("/ali/callback/:purpose", ac.OssOperationCallback)

	// 目前主要有两块业务功能：1. 登录；2. 会话留存
	// 1. 登录功能包含：
	// 	1）登录 验证账号和密码，返回token
	// 	2）注册 采用邮箱注册，并发送验证邮件, 携带token
	//  3）验证邮件 验证token，更改用户认证状态
	//	以上3步都不用token验证身份信息
	//	4) 上传头像
	//	5) 下载头像
	//	6）更改密码

	v1 := g.Group("v1")
	{
		userv1 := v1.Group("users")
		{
			userv1.POST("register", uc.Register)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(middleware.Authn())
			userv1.GET(":name", uc.Get)
			userv1.GET(":name/verify-email", ec.SendVerificationEmail)
		}
		aliv1 := v1.Group("ali")
		{
			aliv1.GET(":purpose/callback", ac.OssOperationCallback)
			aliv1.Use(middleware.Authn())
			aliv1.GET(":purpose/get-policy-token", ac.GetPolicyToken)
		}
	}

	return nil
}
