package store

import (
	"context"
	"jonghong/internal/pkg/model"

	"gorm.io/gorm"
)

// 直接面向数据的层次结构
// 主要是定义增删查改
// 定义接口
type UserStore interface {
	// 传入用户的全部信息以创建客户
	Create(ctx context.Context, user *model.UserM) error
	// 根据用户名查找用户
	Get(ctx context.Context, username string) (*model.UserM, error)
	// 更新用户信息
	Update(ctx context.Context, user *model.UserM) error
}

// 定义实体
type users struct {
	db *gorm.DB
}

// 实体实现接口
func (u *users) Create(ctx context.Context, user *model.UserM) error {

	// 先删除没有验证的用户信息
	u.db.Where("isVerified = ?", false).Delete(&user)
	// 根据结构体自动关联数据库
	return u.db.Create(&user).Error
}
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(&user).Error
}
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var userM model.UserM
	if err := u.db.Where("username = ?", username).First(&userM).Error; err != nil {
		return nil, err
	}
	return &userM, nil
}

func newUsers(db *gorm.DB) UserStore {
	return &users{db: db}
}
