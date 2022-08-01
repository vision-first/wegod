package middlewares

import (
	"github.com/995933447/simpletrace"
	simpletracectx "github.com/995933447/simpletrace/context"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/enum"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/context/http/header"
)

func Trace(ctx *gin.Context) {
	var traceId, spanId string
	traceId = ctx.GetHeader(header.HeaderSimpleTraceId)
	if traceId == "" {
		traceId = simpletrace.NewTraceId()
		spanId = simpletrace.NewSpanId()
	} else {
		spanId = ctx.GetHeader(header.HeaderSimpleTraceSpanId)
	}

	ctx.Set(enum.GinCtxKeySimpleTraceCtx, simpletracectx.New(config.Conf.App.AppName, ctx, traceId, spanId))

	ctx.Next()

	if ctx.Writer.Written() {
		return
	}

	if simpleTraceCtxForAnyType, ok := ctx.Get(enum.GinCtxKeySimpleTraceCtx); ok {
		traceCtx := simpleTraceCtxForAnyType.(*simpletracectx.Context)
		ctx.Header(header.HeaderSimpleTraceId, traceCtx.GetTraceId())
		ctx.Header(header.HeaderSimpleTraceSpanId, traceCtx.GetSpanId())
		ctx.Header(header.HeaderSimpleTraceParentSpanId, traceCtx.GetParentSpanId())
	}
}