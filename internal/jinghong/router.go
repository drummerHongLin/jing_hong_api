package jinghong

import (
	"jonghong/internal/jinghong/controller/v1/ali"
	"jonghong/internal/jinghong/controller/v1/chat"
	"jonghong/internal/jinghong/controller/v1/email"
	"jonghong/internal/jinghong/controller/v1/payment"
	"jonghong/internal/jinghong/controller/v1/user"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	emailservice "jonghong/internal/pkg/emailservice"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// 采用gin框架搭建服务路由
// 路由中可以添加中间件，路由及中间件的调用顺序和添加顺序相关
// 常用的中间件有 验证身份和授权

func initRouter(g *gin.Engine) error {

	// 没有对应功能
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	// 初始化controller
	uc := user.NewUserController(store.S)
	ac := ali.NewAliController(store.S)
	ec := email.NewEmailController(store.S, emailservice.MS)
	cc := chat.NewChatController(store.S)
	pc := payment.NewPaymentController(store.S)

	// 安全策略
	g.Use(middleware.Cors, middleware.NoCache)

	// 登录不需要token验证，采用用户名和密码验证
	// 回调函数也不需要
	g.POST("/login", uc.Login)

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
			userv1.GET(":name/verify", uc.Verify)
			userv1.GET(":name/send-email", ec.SendVerificationEmail)
			userv1.POST(":name/verify-email", ec.VerifyEmail)
			userv1.Use(middleware.Authn())
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.GET("current-user", uc.Get)
			userv1.POST(":name/set-avatar", ac.SetAvatar)
			userv1.GET(":name/get-avatar", ac.GetAvatar)
		}
		chatv1 := v1.Group("chat")
		{ // 这里面的全部要身份认证
			chatv1.Use(middleware.Authn())
			chatv1.POST("create-new-session", cc.CreateNewSession)
			chatv1.POST("create-new-message", cc.CreateNewMessage)
			chatv1.GET("get-messages/:sessionId", cc.GetMessagesBySession)
			chatv1.GET("get-sessions/:chatModel", cc.GetSessionsByModel)
			chatv1.PUT("delete-session/:sessionId", cc.DeleteSession)
			chatv1.PUT("update-message", cc.UpdateMessage)
			chatv1.GET("get-all-messages", cc.GetMessagesBySession)
			chatv1.GET("get-all-sessions", cc.GetSessionsByModel)
		}
		paymentv1 := v1.Group("payment")
		{
			paymentv1.Use(middleware.Authn())
			paymentv1.POST("create-new-payment", pc.CreateNewPaymentRecord)
			paymentv1.POST("update-payment", pc.UpdatePaymentRecord)
			paymentv1.GET("get-payment/:paymentNo", pc.GetPaymentRecordByNo)
			paymentv1.GET("get-payments", pc.GetPaymentRecordsById)
		}
	}

	return nil
}
