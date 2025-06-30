package ali

import (
	"io"
	"jonghong/internal/jinghong/biz"
	"jonghong/internal/jinghong/store"
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/errno"
	"jonghong/internal/pkg/log"

	"net/http"

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
func (ac *AliController) SetAvatar(c *gin.Context) {
	log.C(c).Infow("Set avatar function called")
	// 获取请求中的文件信息
	fileHeader, err := c.FormFile("file")
	if err != nil {
		core.WriteResponse(c, errno.ErrNoFileUpdated, nil)
		return
	}

	username := c.Param("name")

	file, err := fileHeader.Open()
	if err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}

	defer file.Close()

	fileBody, err := io.ReadAll(file)

	if err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}

	fileType := http.DetectContentType(fileBody)

	if !isImageType(fileType) {
		core.WriteResponse(c, errno.ErrFileTypeUnsupport, nil)
		return
	}

	fileKey := username + fileHeader.Filename

	if err := ac.b.AliBiz().PutObject(fileBody, fileKey, http.DetectContentType(fileBody), c); err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}
	// 更新数据库文件地址
	if err := ac.b.UserBiz().SetAvatar(c, username, fileKey); err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, nil)

}

func (ac *AliController) GetAvatar(c *gin.Context) {
	log.C(c).Infow("Get avatar function called")
	username := c.Param("name")
	userM, err := ac.b.UserBiz().Get(c, username)
	if err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}
	if userM.AvatarUrl == "" {
		userM.AvatarUrl = "未登录.png"
	}

	fileBody, fileType, err := ac.b.AliBiz().GetObject(userM.AvatarUrl, c)
	if err != nil {
		core.WriteResponse(c, errno.FileSysError.SetMessage("%s", err.Error()), nil)
		return
	}

	c.Data(http.StatusOK, fileType, fileBody)

}

func isImageType(mimeType string) bool {
	// 常见图片MIME类型
	imageTypes := []string{
		"image/jpeg", // .jpg, .jpeg
		"image/png",  // .png
		"image/gif",  // .gif
	}

	// 方法1：检查是否在预定义图片类型中
	for _, t := range imageTypes {
		if mimeType == t {
			return true
		}
	}

	// 方法2：更宽松的判断 - 检查MIME类型前缀
	// return strings.HasPrefix(mimeType, "image/")

	return false
}
