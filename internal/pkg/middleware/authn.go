package middleware

import (
	"jonghong/internal/pkg/core"
	"jonghong/internal/pkg/known"
	"jonghong/pkg/token"

	"github.com/gin-gonic/gin"
)

// 验证token
func Authn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := token.ParseRequest(ctx)
		if err != nil {
			core.WriteResponse(ctx, err, nil)
			// 停止请求的下一步处理
			ctx.Abort()
			return
		}
		ctx.Set(known.XUsernameKey, username)
		ctx.Next()
	}
}
