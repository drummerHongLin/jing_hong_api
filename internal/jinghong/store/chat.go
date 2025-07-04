package store

import (
	"context"
	"jonghong/internal/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChatStore interface {
	CreateSession(ctx context.Context, session []*model.SessionM) error
	GetSessions(ctx context.Context, chatModel string, userId uint) ([]model.SessionM, error)
	GetSession(ctx context.Context, id string, userId uint) (model.SessionM, error)
	DeleteSession(ctx context.Context, session *model.SessionM) error
	CreateMessage(ctx context.Context, message []*model.MessageM) error
	// 用外键查
	GetMessages(ctx context.Context, session *model.SessionM) ([]model.MessageM, error)
	UpdateSession(ctx context.Context, session *model.SessionM) error
	UpdateMessage(ctx context.Context, message *model.MessageM) error
	GetAllMessages(ctx context.Context, userId uint) ([]model.MessageM, error)
	GetAllSessions(ctx context.Context, userId uint) ([]model.SessionM, error)
}

type chat struct {
	db *gorm.DB
}

func newChat(db *gorm.DB) ChatStore {
	return &chat{db: db}
}

func (c *chat) CreateSession(ctx context.Context, session []*model.SessionM) error {
	// 这样前端就不用查询全部的信息，浪费流量
	return c.db.Clauses(clause.OnConflict{DoNothing: true}).Create(session).Error
}

func (c *chat) DeleteSession(ctx context.Context, session *model.SessionM) error {
	return c.db.Delete(&session).Error
}

func (c *chat) GetSessions(ctx context.Context, chatModel string, userId uint) ([]model.SessionM, error) {
	var sessions []model.SessionM
	err := c.db.Where("chatModel = ?", chatModel).Where("userId = ?", userId).Order("createTime desc").Find(&sessions).Error
	return sessions, err
}

func (c *chat) GetAllSessions(ctx context.Context, userId uint) ([]model.SessionM, error) {
	var sessions []model.SessionM
	err := c.db.Where("userId = ?", userId).Order("createTime desc").Find(&sessions).Error
	return sessions, err
}

func (c *chat) GetSession(ctx context.Context, id string, userId uint) (model.SessionM, error) {
	var session model.SessionM
	err := c.db.Where(&model.SessionM{ID: id, UserId: userId}).First(&session).Error
	return session, err

}

func (c *chat) CreateMessage(ctx context.Context, message []*model.MessageM) error {
	return c.db.Clauses(clause.OnConflict{DoNothing: true}).Create(message).Error
}

func (c *chat) GetMessages(ctx context.Context, session *model.SessionM) ([]model.MessageM, error) {
	var messages []model.MessageM
	err := c.db.Where("sessionId = ?", session.ID).Order("sendTime").Find(&messages).Error
	return messages, err
}

func (c *chat) GetAllMessages(ctx context.Context, userId uint) ([]model.MessageM, error) {
	var sessions []model.SessionM
	if err := c.db.Where("userId = ?", userId).Order("createTime").Find(&sessions).Error; err != nil {
		return nil, err
	}

	var sId []string

	for _, s := range sessions {
		sId = append(sId, s.ID)
	}

	var messages []model.MessageM
	err := c.db.Where("sessionId in ?", sId).Order("sendTime desc").Find(&messages).Error
	return messages, err
}

func (c *chat) UpdateSession(ctx context.Context, session *model.SessionM) error {
	return c.db.Save(&session).Error
}

func (c *chat) UpdateMessage(ctx context.Context, message *model.MessageM) error {
	return c.db.Save(&message).Error
}
