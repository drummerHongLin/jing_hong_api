package email

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	emailservice "jonghong/internal/pkg/emailservice"
	"jonghong/internal/pkg/errno"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/gin-gonic/gin"
)

type EmailController struct {
	b  biz.IBiz
	ms emailservice.MailService
}

func NewEmailController(ds store.IStore, ms emailservice.MailService) *EmailController {
	return &EmailController{
		b:  biz.NewBiz(ds),
		ms: ms,
	}
}

func (ec *EmailController) SendVerificationEmail(c *gin.Context) {
	username := c.Param("name")
	if username == "" {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	_, err := ec.b.UserBiz().Get(c, username)

	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if err := ec.b.EmailBiz(ec.ms).SendVerificationEmail(c, username); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "发送邮件成功!")

}

func (ec *EmailController) VerifyEmail(c *gin.Context) {
	username := c.Param("name")

	var r v1.EmailVerifingRequest

	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	rs, err := ec.b.EmailBiz(ec.ms).VerifyEmail(c, username, r.Code)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	// 如果验证成功返回token
	core.WriteResponse(c, nil, rs)
	//c.HTML(http.StatusOK, "SuccessVerified.html", gin.H{
	//	"title":   username,
	//	"welcome": "您已认证成功，欢迎使用JingHong！",
	//})

}
