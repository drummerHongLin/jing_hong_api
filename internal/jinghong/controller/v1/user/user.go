package user

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	v1 "jonghong/pkg/api/jinghong/v1"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 解析http请求信息，验证请求消息是否有效

type UserController struct {
	b biz.IBiz
}

func NewUserController(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}

// 定义方法

func (uc *UserController) Login(c *gin.Context) {
	log.C(c).Infow("Login function called")

	var r v1.LoginRequest

	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	if err := validator.Struct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("%s", err.Error()), nil)
		return
	}

	resp, err := uc.b.UserBiz().Login(c, &r)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)

}

func (uc *UserController) Register(c *gin.Context) {
	log.C(c).Infow("Register function called")

	var r v1.CreateUserRequest

	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	if err := validator.Struct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("%s", err.Error()), nil)
		return
	}

	err := uc.b.UserBiz().Register(c, &r)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)

}

func (uc *UserController) Get(c *gin.Context) {
	log.C(c).Infow("Find user by username called")

	// 从中间件设定中获取账号信息
	user, err := uc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, user)
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("Change password function called")

	var r v1.ChangePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	if err := validator.Struct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("%s", err.Error()), nil)
		return
	}

	// 从中间件设定中获取账号信息
	if err := uc.b.UserBiz().ChangePassword(c, c.GetString(known.XUsernameKey), &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "修改密码成功!")

}

func (uc *UserController) Verify(c *gin.Context) {
	log.C(c).Infow("Verify user function called")

	username := c.Param("name")

	userM, err := uc.b.UserBiz().Get(c, username)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	if !userM.IsVerified {
		core.WriteResponse(c, errno.ErrUserNotVerified, nil)
		return
	}

	core.WriteResponse(c, nil, nil)

}
