package chat

import (
	"context"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/model"
	v1 "jonghong/pkg/api/jinghong/v1"
	"regexp"

	"github.com/jinzhu/copier"
)

type ChatBiz interface {
	CreateNewSession(ctx context.Context, r []v1.NewSessionRequest, userId uint) error
	CreateNewMessage(ctx context.Context, r []v1.NewMessageRequest) error
	GetMessagesBySession(ctx context.Context, sessionId string, userId uint) (*v1.GetMessagesResponse, error)
	GetSessionsByModel(ctx context.Context, chatModel string, userId uint) (*v1.GetSessionsResponse, error)
	UpdateMessage(ctx context.Context, r *v1.NewMessageRequest) error
	DeleteSession(ctx context.Context, sessionId string, userId uint) error
	GetAllSessions(ctx context.Context, userId uint) (*v1.GetSessionsResponse, error)
	GetAllMessages(ctx context.Context, userId uint) (*v1.GetMessagesResponse, error)
}

type chatBiz struct {
	cs store.ChatStore
}

func NewChatBiz(cs store.ChatStore) ChatBiz {
	return &chatBiz{
		cs: cs,
	}
}

func (cb *chatBiz) CreateNewSession(ctx context.Context, r []v1.NewSessionRequest, userId uint) error {
	var sessions []*model.SessionM
	_ = copier.Copy(&sessions, r)
	for _, v := range sessions {
		v.UserId = userId
	}
	err := cb.cs.CreateSession(ctx, sessions)
	if err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'id'", err.Error()); match {
			return errno.ErrSessionAlreadyExist
		}
		return err
	}
	return nil
}

func (cb *chatBiz) CreateNewMessage(ctx context.Context, r []v1.NewMessageRequest) error {
	var messages []*model.MessageM
	_ = copier.Copy(&messages, r)
	err := cb.cs.CreateMessage(ctx, messages)
	if err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'id'", err.Error()); match {
			return errno.ErrSessionAlreadyExist.SetMessage("消息ID重复!")
		}
		return err
	}
	return nil
}

func (cb *chatBiz) GetMessagesBySession(ctx context.Context, sessionId string, userId uint) (*v1.GetMessagesResponse, error) {

	session, err := cb.cs.GetSession(ctx, sessionId, userId)

	if err != nil {
		return nil, err
	}

	messages, err := cb.cs.GetMessages(ctx, &session)

	if err != nil {
		return nil, err
	}

	var returnMessages []v1.Message

	// 这一步不一定正确，需要单元测试来着
	_ = copier.Copy(&returnMessages, messages)

	return &v1.GetMessagesResponse{Messages: returnMessages}, nil

}

func (cb *chatBiz) GetSessionsByModel(ctx context.Context, chatModel string, userId uint) (*v1.GetSessionsResponse, error) {

	session, err := cb.cs.GetSessions(ctx, chatModel, userId)

	if err != nil {
		return nil, err
	}

	var returnSessions []v1.Session

	_ = copier.Copy(&returnSessions, session)

	return &v1.GetSessionsResponse{Sessions: returnSessions}, nil
}

func (cb *chatBiz) UpdateMessage(ctx context.Context, r *v1.NewMessageRequest) error {
	var message model.MessageM
	_ = copier.Copy(&message, r)

	if err := cb.cs.UpdateMessage(ctx, &message); err != nil {
		return err
	}

	return nil

}
func (cb *chatBiz) DeleteSession(ctx context.Context, sessionId string, userId uint) error {
	session, err := cb.cs.GetSession(ctx, sessionId, userId)

	if err != nil {
		return err
	}

	if err := cb.cs.DeleteSession(ctx, &session); err != nil {
		return err
	}
	return nil

}

func (cb *chatBiz) GetAllSessions(ctx context.Context, userId uint) (*v1.GetSessionsResponse, error) {
	session, err := cb.cs.GetAllSessions(ctx, userId)

	if err != nil {
		return nil, err
	}

	var returnSessions []v1.Session

	_ = copier.Copy(&returnSessions, session)

	return &v1.GetSessionsResponse{Sessions: returnSessions}, nil
}

func (cb *chatBiz) GetAllMessages(ctx context.Context, userId uint) (*v1.GetMessagesResponse, error) {

	messages, err := cb.cs.GetAllMessages(ctx, userId)

	if err != nil {
		return nil, err
	}

	var returnMessages []v1.Message

	// 这一步不一定正确，需要单元测试来着
	_ = copier.Copy(&returnMessages, messages)

	return &v1.GetMessagesResponse{Messages: returnMessages}, nil

}
