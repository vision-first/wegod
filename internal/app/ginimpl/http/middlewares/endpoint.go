package middlewares

import (
	"github.com/gin-gonic/gin"
	"wegod/internal/app/ginimpl/http/response"
)

func RspFlusher(ctx *gin.Context) {
	ctx.Next()
	response.Flush(ctx)
}
