package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
)

func RspFlusher(ctx *gin.Context) {
	ctx.Next()
	response.Flush(ctx)
}
