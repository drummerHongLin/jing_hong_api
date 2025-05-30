package email

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	emailservice "jonghong/internal/pkg/emailservicee"
	"jonghong/internal/pkg/errno"
	"jonghong/pkg/token"
	"net/http"

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
	if err := ec.b.EmailBiz(ec.ms).SendVerificationEmail(c, username); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "发送邮件成功!")

}

func (ec *EmailController) VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		core.WriteResponse(c, errno.ErrTokenInvalid, nil)
		return
	}

	username, err := token.Parse(tokenString)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := ec.b.EmailBiz(ec.ms).VerifyEmail(c, username); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	c.HTML(http.StatusOK, "SuccessVerified.html", gin.H{
		"title":   username,
		"welcome": "您已认证成功，欢迎使用JingHong！",
	})

}
