package user

import (
	"context"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/model"
	v1 "jonghong/pkg/api/jinghong/v1"
	"jonghong/pkg/auth"
	"jonghong/pkg/token"
	"regexp"
	"time"

	"github.com/jinzhu/copier"
)

// 进一步封装数据功能
// 实现从接口到数据库的对应

type UserBiz interface {
	Register(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	SetAvatar(ctx context.Context, username string, avatarUrl string) error
	Get(ctx context.Context, username string) (*model.UserM, error)
}

// userBiz不对外暴露，统一在biz中

type userBiz struct {
	us store.UserStore
}

func NewUserBiz(us store.UserStore) *userBiz {
	return &userBiz{us: us}
}

// 实现接口逻辑

func (ub *userBiz) Register(ctx context.Context, r *v1.CreateUserRequest) error {
	//创建一个UserM对象，使请求对象转化为UserM对象
	var userM model.UserM
	// 这个包实现根据相同字段进行转换
	_ = copier.Copy(&userM, r)
	// 在数据库中创建记录
	err := ub.us.Create(ctx, &userM)
	// 判断用户名是否已存在
	if err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil

}

func (ub *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 根据用户名从数据库中获取用户信息
	userM, err := ub.us.Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}
	// 比较保存的密码和输入的密码
	if err := auth.Compare(userM.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}
	// 签发token
	expiredAt := time.Now().Add(7 * 24 * time.Hour).Unix()
	t, err := token.Sign(r.Username, expiredAt)

	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t, ExpiredAt: expiredAt}, nil

}

func (ub *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := ub.us.Get(ctx, username)
	if err != nil {
		return err
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)

	if err := ub.us.Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

func (ub *userBiz) SetAvatar(ctx context.Context, username string, avatarUrl string) error {
	userM, err := ub.us.Get(ctx, username)
	if err != nil {
		return err
	}
	// 在路由中做身份拦截，这里不需要验证密码
	userM.AvatarUrl = avatarUrl
	if err := ub.us.Update(ctx, userM); err != nil {
		return err
	}
	return nil
}

func (ub *userBiz) Get(ctx context.Context, username string) (*model.UserM, error) {
	userM, err := ub.us.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return userM, nil
}
