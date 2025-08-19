package player

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"
	"jonghong/internal/pkg/log"
	v1 "jonghong/pkg/api/jinghong/v1"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerController struct {
	b biz.IBiz
}

func NewPlayerController(ds store.IStore) PlayerController {
	return PlayerController{b: biz.NewBiz(ds)}
}

func (pc *PlayerController) CreatePlayConfig(c *gin.Context) {
	log.C(c).Infow("Create new player config function called")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	var r v1.PlayConfig
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	err = pc.b.PlayerBiz().CreatePlayConfig(c, &r, user.ID)

	core.WriteResponse(c, err, nil)
}
func (pc *PlayerController) UpdatePlayConfig(c *gin.Context) {
	log.C(c).Infow("Update player config function called")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	var r v1.PlayConfig
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	err = pc.b.PlayerBiz().UpdatePlayConfig(c, &r, user.ID)

	core.WriteResponse(c, err, nil)
}
func (pc *PlayerController) GetPlayConfigs(c *gin.Context) {
	log.C(c).Infow("Get PlayConfigs by Id function called")
	// 解析参数变量
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("非法offset参数，参数必须为整数"), nil)
		return
	}
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage("非法limit参数，参数必须为整数"), nil)
		return
	}
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	record, err := pc.b.PlayerBiz().GetPlayConfigs(c, offset, limit, user.ID)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, record)
}
func (pc *PlayerController) InsertPlayConfigs(c *gin.Context) {
	log.C(c).Infow("Imsert Play Config  function called")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	var r []v1.PlayConfigWithTime
	if err := c.ShouldBindBodyWithJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	err = pc.b.PlayerBiz().InsertPlayConfigs(c, r, user.ID)
	core.WriteResponse(c, err, nil)
}
func (pc *PlayerController) DeletePlayConfig(c *gin.Context) {
	log.C(c).Infow("Delete Play Config  function called")
	playConfigNo := c.Param("playConfigNo")
	// 先通过token查找user
	user, err := pc.b.UserBiz().Get(c, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	err = pc.b.PlayerBiz().DeletePlayConfig(c, user.ID, playConfigNo)
	core.WriteResponse(c, err, nil)
}
