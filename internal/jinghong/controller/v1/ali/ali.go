package ali

import (
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/known"

	"github.com/gin-gonic/gin"
)

type AliController struct {
	b biz.IBiz
}

func NewAliController(ds store.IStore) *AliController {
	return &AliController{b: biz.NewBiz(ds)}
}

// 目前只需要签名信息
// 前端在获取签名信息后通过oss接口，和ali交互

func (ac *AliController) GetPolicyToken(c *gin.Context) {
	res, err := ac.b.AliBiz().GetPolicyToken(c.GetString(known.XUsernameKey), c.Param("purpose"))
	if err != nil {
		core.WriteResponse(c, errno.InternalServerError.SetMessage("%s", err.Error()), nil)
		return
	}
	core.WriteResponse(c, nil, res)
}

// 需要设置回调
// 回调获取文件信息

func (ac *AliController) OssOperationCallback(c *gin.Context) {

	// 3个参数缺一不可
	var purpose = c.Param("purpose")
	var filename = c.Query("filename")
	var username = c.Query("username")

	if purpose == "" || filename == "" || username == "" {
		core.WriteResponse(c, errno.ErrInvalidParameter, nil)
		return
	}

	// 设置头像先占一个坑位
	if purpose == "set-avatar" {
		avatarUrl := ac.b.AliBiz().GetOssHost() + "/" + filename
		if err := ac.b.UserBiz().SetAvatar(c, username, avatarUrl); err != nil {
			core.WriteResponse(c, err, nil)
		}
	}

	core.WriteResponse(c, nil, "回调成功")

}
