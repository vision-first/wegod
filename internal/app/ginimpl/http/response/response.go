package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Content interface{}
	StatusCode int
	ContentFmt RspContentFmt
}

type JsonResponseContent struct {
	ErrCode int
	Data interface{}
	ErrMsg string
}

type RspContentFmt int

const (
	CtxKeyRspData = "http-rsp-data"
)

const (
	RspContentFmtJson RspContentFmt = iota
)

func EndSuccessfulJson(ctx *gin.Context, data interface{}) {
	ctx.Set(CtxKeyRspData, &Response{
		StatusCode: http.StatusOK,
		Content: &JsonResponseContent{
			Data: data,
		},
		ContentFmt: RspContentFmtJson,
	})
}

func FlushByJson(ctx *gin.Context) {
	if rspForInterface, ok := ctx.Get(CtxKeyRspData); ok {
		rsp := rspForInterface.(*Response)
		ctx.JSON(rsp.StatusCode, rsp.Content)
	}
}

func Flush(ctx *gin.Context) {
	if ctx.Writer.Written() {
		return
	}

	if rspForInterface, ok := ctx.Get(CtxKeyRspData); ok {
		switch rspForInterface.(*Response).ContentFmt {
		case RspContentFmtJson:
			FlushByJson(ctx)
		default:
			FlushByJson(ctx)
		}
		return
	}

	FlushByJson(ctx)
}
