package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/middlewares"
	"github.com/vision-first/wegod/internal/app/ginimpl/http/response"
	"github.com/vision-first/wegod/internal/pkg/api/apis"
	"github.com/vision-first/wegod/internal/pkg/boot"
)

type HttpRouterProvider struct {
	srv *gin.Engine
}

func NewHttpRouterProvider(srv *gin.Engine) *HttpRouterProvider {
	return &HttpRouterProvider{
		srv: srv,
	}
}

var _ boot.ServiceProvider = (*HttpRouterProvider)(nil)

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

	buddhaApi := apis.NewBuddha(logger)
	postApi := apis.NewPost(logger)
	authApi := apis.NewAuth(logger)

	apiDispatcher := ginimpl.NewApiDispatcher()
	publicGroup := p.srv.Group("/pub")
	publicGroup.GET("/buddhas", apiDispatcher.MakeDispatchFunc(buddhaApi.PageBuddha))
	publicGroup.GET("/posts", apiDispatcher.MakeDispatchFunc(postApi.PagePosts))
	publicGroup.GET("/post/categories", apiDispatcher.MakeDispatchFunc(postApi.PageCategories))
	publicGroup.GET("/post/post", apiDispatcher.MakeDispatchFunc(postApi.GetPost))
	publicGroup.POST("/register", apiDispatcher.MakeDispatchFunc(authApi.Register))
	publicGroup.POST("/login", apiDispatcher.MakeDispatchFunc(authApi.Login))

	return nil
}
