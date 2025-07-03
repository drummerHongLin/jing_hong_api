package biz

import (
	"jonghong/internal/jinghong/biz/ali"
	"jonghong/internal/jinghong/biz/chat"
	"jonghong/internal/jinghong/biz/email"
	"jonghong/internal/jinghong/biz/user"
	"jonghong/internal/jinghong/store"
	emailservice "jonghong/internal/pkg/emailservice"
)

// biz层的统一入口

type IBiz interface {
	UserBiz() user.UserBiz
	AliBiz() ali.AliBiz
	ChatBiz() chat.ChatBiz
	EmailBiz(ms emailservice.MailService) email.EmailBiz
}

type biz struct {
	ds store.IStore
}

func (b *biz) UserBiz() user.UserBiz {
	return user.NewUserBiz(b.ds.Users())
}

func (b *biz) ChatBiz() chat.ChatBiz {
	return chat.NewChatBiz(b.ds.Chat())
}

func (b *biz) AliBiz() ali.AliBiz {
	return ali.NewAliBiz()
}

func (b *biz) EmailBiz(ms emailservice.MailService) email.EmailBiz {
	return email.NewEmailBiz(ms, b.ds.Users())
}

func NewBiz(ds store.IStore) IBiz {
	return &biz{ds: ds}
}
