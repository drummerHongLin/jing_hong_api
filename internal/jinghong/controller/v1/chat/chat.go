package chat

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	b biz.IBiz
}

func NewChatController(ds store.IStore) ChatController {
	return ChatController{b: biz.NewBiz(ds)}
}

func (cc *ChatController) CreateNewSession(c *gin.Context) {
	log.C(c).Infow("Create new session function called")

	var r []v1.NewSessionRequest

	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	// 获取用户名
	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := cc.b.ChatBiz().CreateNewSession(c, r, uint(user.ID)); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)

}

func (cc *ChatController) CreateNewMessage(c *gin.Context) {
	log.C(c).Infow("Create new message function called")

	var r []v1.NewMessageRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if err := cc.b.ChatBiz().CreateNewMessage(c, r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)

}

func (cc *ChatController) GetMessagesBySession(c *gin.Context) {
	log.C(c).Infow("Get messages function called")

	sessionId := c.Param("sessionId")

	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	messages, err := cc.b.ChatBiz().GetMessagesBySession(c, sessionId, uint(user.ID))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, messages)

}

func (cc *ChatController) GetSessionsByModel(c *gin.Context) {
	log.C(c).Infow("Get sessions function called")
	chatModel := c.Param("chatModel")
	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	sessions, err := cc.b.ChatBiz().GetSessionsByModel(c, chatModel, uint(user.ID))

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, sessions)

}

func (cc *ChatController) UpdateMessage(c *gin.Context) {
	log.C(c).Infow("Create new message function called")

	var r v1.NewMessageRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if err := cc.b.ChatBiz().UpdateMessage(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)

}

func (cc *ChatController) DeleteSession(c *gin.Context) {
	log.C(c).Infow("Delete session function called")

	sessionId := c.Param("sessionId")

	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := cc.b.ChatBiz().DeleteSession(c, sessionId, uint(user.ID)); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}

func (cc *ChatController) GetAllSessions(c *gin.Context) {
	log.C(c).Infow("Get all sessions function called")
	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	sessions, err := cc.b.ChatBiz().GetAllSessions(c, uint(user.ID))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, sessions)

}

func (cc *ChatController) GetAllMessages(c *gin.Context) {
	log.C(c).Infow("Get all messages function called")
	user, err := cc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	messages, err := cc.b.ChatBiz().GetAllMessages(c, uint(user.ID))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, messages)

}
