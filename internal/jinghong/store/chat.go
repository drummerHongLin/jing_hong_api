package store

import (
	"context"
	"jonghong/internal/pkg/model"

	"gorm.io/gorm"
)

type ChatStore interface {
	CreateSession(ctx context.Context, session *model.SessionM) error
	GetSessions(ctx context.Context, chatModel string, userId uint) ([]model.SessionM, error)
	GetSession(ctx context.Context, id uint, userId uint) (model.SessionM, error)
	DeleteSession(ctx context.Context, session *model.SessionM) error
	CreateMessage(ctx context.Context, message *model.MessageM) error
	// 用外键查
	GetMessages(ctx context.Context, session *model.SessionM) ([]model.MessageM, error)
	UpdateSession(ctx context.Context, session *model.SessionM) error
	UpdateMessage(ctx context.Context, message *model.MessageM) error
}

type chat struct {
	db *gorm.DB
}

func newChat(db *gorm.DB) ChatStore {
	return &chat{db: db}
}

func (c *chat) CreateSession(ctx context.Context, session *model.SessionM) error {
	return c.db.Create(&session).Error
}

func (c *chat) DeleteSession(ctx context.Context, session *model.SessionM) error {
	return c.db.Delete(&session).Error
}

func (c *chat) GetSessions(ctx context.Context, chatModel string, userId uint) ([]model.SessionM, error) {
	var sessions []model.SessionM
	err := c.db.Where("chatModel = ?", chatModel).Where("userId = ?", userId).Find(&sessions).Error
	return sessions, err
}

func (c *chat) GetSession(ctx context.Context, id uint, userId uint) (model.SessionM, error) {
	var session model.SessionM
	err := c.db.Where(&model.SessionM{ID: id, UserId: userId}).First(&session).Error
	return session, err

}

func (c *chat) CreateMessage(ctx context.Context, message *model.MessageM) error {
	return c.db.Create(&message).Error
}

func (c *chat) GetMessages(ctx context.Context, session *model.SessionM) ([]model.MessageM, error) {
	var messages []model.MessageM
	err := c.db.Where("sessionId = ?", session.ID).Find(&messages).Error
	return messages, err
}

func (c *chat) UpdateSession(ctx context.Context, session *model.SessionM) error {
	return c.db.Save(&session).Error
}

func (c *chat) UpdateMessage(ctx context.Context, message *model.MessageM) error {
	return c.db.Save(&message).Error
}
