package email

import (
	"bytes"
	"context"
	"html/template"
	"jonghong/internal/jinghong/store"
	emailservice "jonghong/internal/pkg/emailservice"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	v1 "jonghong/pkg/api/jinghong/v1"
	"jonghong/pkg/token"
	"path/filepath"
	"time"
)

// 定义email事务接口

type EmailBiz interface {
	SendVerificationEmail(ctx context.Context, username string) error
	VerifyEmail(ctx context.Context, username string, code string) (*v1.EmailVerifiedResponse, error)
}

type emailBiz struct {
	ms emailservice.MailService
	us store.UserStore
}

func NewEmailBiz(ms emailservice.MailService, us store.UserStore) EmailBiz {
	return &emailBiz{ms: ms, us: us}
}

// 1. 发送邮件
// 发送邮件调用的身份验证采用token拦截
// 携带username参数，生成token后续验证中提取出来
// token的时限为1小时

func (eb *emailBiz) SendVerificationEmail(ctx context.Context, username string) error {

	userM, err := eb.us.Get(ctx, username)
	if err != nil {
		return errno.ErrUserNotFound
	}

	// 生成验证码
	code, err := eb.ms.GenerateCode(username)
	if err != nil {
		return err
	}

	temlPath := filepath.Join(known.HomeDir, "/static/html/EmailVerification.html")

	// 组装邮件body
	// 采用相对地址， 静态文件夹放在可执行文件的同一级
	t := template.Must(template.ParseFiles(temlPath))

	var buf bytes.Buffer

	if err := t.Execute(&buf, map[string]any{
		"username":    username,
		"code":        code,
		"currentDate": time.Now().Format("2006-01-02"),
	}); err != nil {
		return err
	}

	msg := emailservice.MailMessage{
		From:    "",
		To:      []string{userM.Email},
		Subject: "JingHong - 邮件验证",
		Body:    buf.String(),
	}

	if err := eb.ms.SendEmailAsync(&msg); err != nil {
		return err
	}

	return nil
}

func (eb *emailBiz) VerifyEmail(ctx context.Context, username string, code string) (*v1.EmailVerifiedResponse, error) {

	userM, err := eb.us.Get(ctx, username)

	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	rst, err := eb.ms.VerifyCode(username, code)

	if err != nil {
		return nil, err
	}
	if !rst {
		return nil, errno.ErrCodeNotExist
	}

	userM.IsVerified = true
	if err := eb.us.Update(ctx, userM); err != nil {
		return nil, errno.InternalServerError.SetMessage("数据更新失败，请重试")
	}

	// 签发token
	t, err := token.Sign(username, time.Now().Add(7*24*time.Hour).Unix())

	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.EmailVerifiedResponse{Token: t}, nil
}
