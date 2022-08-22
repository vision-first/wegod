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
	userApi := apis.NewUser(logger)
	shopApi := apis.NewShop(logger)
	prayPropApi := apis.NewPrayProp(logger)
	wishPropApi := apis.NewWorshipProp(logger)
	musicApi := apis.NewMusic(logger)
	donationApi := apis.NewDonation(logger)

	apiDispatcher := ginimpl.NewApiDispatcher()

	publicGroup := p.srv.Group("/pub")
	privateGroup := p.srv.Group("/private").Use(middlewares.Auth)

	publicGroup.GET("/buddhas", apiDispatcher.MakeDispatchFunc(buddhaApi.PageBuddha))
	publicGroup.GET("/posts", apiDispatcher.MakeDispatchFunc(postApi.PagePosts))
	publicGroup.GET("/post/categories", apiDispatcher.MakeDispatchFunc(postApi.PageCategories))
	publicGroup.GET("/post/post", apiDispatcher.MakeDispatchFunc(postApi.GetPost))
	publicGroup.POST("/register", apiDispatcher.MakeDispatchFunc(authApi.Register))
	publicGroup.POST("/login", apiDispatcher.MakeDispatchFunc(authApi.Login))
	publicGroup.POST("/register/verify_code", apiDispatcher.MakeDispatchFunc(authApi.SendVerifyCodeForRegister))
	publicGroup.POST("/login/verify_code", apiDispatcher.MakeDispatchFunc(authApi.SendVerifyCodeForLogin))
	publicGroup.GET("/shop/products", apiDispatcher.MakeDispatchFunc(shopApi.PageProducts))
	publicGroup.GET("/shop/product", apiDispatcher.MakeDispatchFunc(shopApi.GetProduct))
	publicGroup.GET("/buddha/rent/packages", apiDispatcher.MakeDispatchFunc(buddhaApi.PageBuddhaRentPackages))
	publicGroup.GET("/buddha/pray/props", apiDispatcher.MakeDispatchFunc(prayPropApi.PagePrayProps))
	publicGroup.GET("/buddha/worship/props", apiDispatcher.MakeDispatchFunc(wishPropApi.PageWorshipProps))
	publicGroup.GET("/musics", apiDispatcher.MakeDispatchFunc(musicApi.PageMusics))
	publicGroup.GET("/shop/product/categories", apiDispatcher.MakeDispatchFunc(shopApi.PageProductCategories))

	privateGroup.POST("/user", apiDispatcher.MakeDispatchFunc(userApi.SetUserInfo))
	privateGroup.GET("/", apiDispatcher.MakeDispatchFunc(userApi.GetUserInfo))
	privateGroup.POST("/buddha/watch", apiDispatcher.MakeDispatchFunc(buddhaApi.WatchBuddha))
	privateGroup.POST("/buddha/unwatch", apiDispatcher.MakeDispatchFunc(buddhaApi.UnwatchBuddha))
	privateGroup.GET("/user/buddha", apiDispatcher.MakeDispatchFunc(buddhaApi.PageUserWatchedBuddhas))
	privateGroup.POST("/buddha/rent/order", apiDispatcher.MakeDispatchFunc(buddhaApi.CreateBuddhaRentOrder))
	privateGroup.POST("/buddha/worship", apiDispatcher.MakeDispatchFunc(buddhaApi.WorshipToBuddha))
	privateGroup.POST("/buddha/pray", apiDispatcher.MakeDispatchFunc(buddhaApi.PrayToBuddha))
	privateGroup.POST("buddha/pray/prop/order", apiDispatcher.MakeDispatchFunc(prayPropApi.CreatePrayPropOrder))
	privateGroup.POST("buddha/worship/prop/order", apiDispatcher.MakeDispatchFunc(wishPropApi.CreateWorshipPropOrder))
	privateGroup.GET("shop/order", apiDispatcher.MakeDispatchFunc(shopApi.GetProduct))
	privateGroup.GET("/shop/orders", apiDispatcher.MakeDispatchFunc(shopApi.PageOrders))
	privateGroup.POST("shop/order", apiDispatcher.MakeDispatchFunc(shopApi.CreateOrder))
	privateGroup.POST("/donation", apiDispatcher.MakeDispatchFunc(donationApi.CreateDonationOrder))
	privateGroup.POST("/donation/ranks", apiDispatcher.MakeDispatchFunc(donationApi.PageDonationRanks))

	return nil
}
