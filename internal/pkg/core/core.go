package core

import (
	"jonghong/internal/pkg/errno"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 这里统一返回的接口
// 如果是错误返回，返回体中携带，内部错误码和错误消息
// 结构体中定义序列化的准则
type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		// 从错误中解析出错误码和信息
		hcode, code, message := errno.DecodeError(err)
		// 返回JSON类型
		c.JSON(hcode, ErrResponse{code, message})
		return
	}
	c.JSON(http.StatusOK, data)
}
