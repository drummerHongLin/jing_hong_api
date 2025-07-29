package model

import (
	"jonghong/pkg/auth"
	"time"

	"gorm.io/gorm"
)

// 定义用户模块数据库表的结构
// 用户的增删改查主要是对用户信息的调整
// token信息直接用jwt进行签发和解码，不需要存数据库
type UserM struct {
	ID         int    `gorm:"column:id;" json:"id"`
	Username   string `gorm:"column:username;not null;unique" json:"username"`
	Password   string `gorm:"column:password;not null" json:"password"`
	Nickname   string `gorm:"column:nickname" json:"nickname"`
	Email      string `gorm:"column:email;not null;default:false" json:"email"`
	IsVerified bool   `gorm:"column:isVerified;not null" json:"isVerified"`
	AvatarUrl  string `gorm:"column:avatarUrl;" json:"avatarUrl"`
	// CreatedAt UpdatedAt 是保留关键字
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

// 实现接口规范
func (u *UserM) TableName() string {
	return "user"
}

// 在创建数据记录前，将明文密码转化成密码
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}
