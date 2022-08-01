package providers

import (
	"github.com/gin-gonic/gin"
	"wegod/internal/app/ginimpl/facades"
	"wegod/internal/app/ginimpl/http/middlewares"
	"wegod/internal/app/ginimpl/http/response"
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

	p.srv.GET("/hello", func(ctx *gin.Context) {
		facades.MustLogger().Debug(ctx, "this is a gin log")
		response.EndSuccessfulJson(ctx, gin.H{
			"hello": "world",
		})
	})

	return nil
}
