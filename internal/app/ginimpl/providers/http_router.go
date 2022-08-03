package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/middlewares"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
	"github.com/vision-first/wegod/internal/pkg/api/apis"
)

type HttpRouterProvider struct {
	srv *gin.Engine
}

func NewHttpRouterProvider(srv *gin.Engine) *HttpRouterProvider {
	return &HttpRouterProvider{
		srv: srv,
	}
}

func (p *HttpRouterProvider) Boot() error {
	p.srv.Use(gin.Recovery())
	p.srv.Use(middlewares.Cors)
	p.srv.Use(middlewares.Trace)
	p.srv.Use(middlewares.RspFlusher)

	logger := facades.MustLogger()
	p.srv.GET("/hello", func(ctx *gin.Context) {
		logger.Debug(ctx, "this is a gin log")
		response.EndSuccessfulJson(ctx, gin.H{
			"hello": "world",
		})
	})

	apiDispatcher := ginimpl.NewApiDispatcher()
	buddhaApi := apis.NewBuddha(logger)
	p.srv.POST("/buddhas", apiDispatcher.MakeDispatchHandlerFunc(buddhaApi.PageBuddha))

	return nil
}
