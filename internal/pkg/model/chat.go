package model

import (
	"time"
)

type MessageM struct {
	ID             string   `gorm:"column:id;" json:"id"`
	MID            uint     `gorm:"column:mid;not null;" json:"mid"`
	Content        string   `gorm:"column:content;not null;" json:"content"`
	Role           string   `gorm:"column:role;not null" json:"role"`
	State          string   `gorm:"column:state;not null" json:"state"`
	ShowingContent string   `gorm:"column:showingContent;not null;" json:"showingContent"`
	SendTime       string   `gorm:"column:sendTime;not null" json:"sendTime"`
	SessionId      string   `gorm:"column:sessionId;not null" json:"sessionId"`
	Session        SessionM `gorm:"foreignKey:SessionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// CreatedAt UpdatedAt 是保留关键字
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// 实现接口规范
func (m *MessageM) TableName() string {
	return "message"
}

type SessionM struct {
	ID         string `gorm:"column:id;" json:"id"`
	Title      string `gorm:"column:title;not null;" json:"title"`
	ChatModel  string `gorm:"column:chatModel;not null;" json:"chatModel"`
	CreateTime string `gorm:"column:CreateTime;not null" json:"sendTime"`
	UserId     uint   `gorm:"column:userId; not null" json:"userId"`
	User       UserM  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// CreatedAt UpdatedAt 是保留关键字
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// 实现接口规范
func (s *SessionM) TableName() string {
	return "session"
}
