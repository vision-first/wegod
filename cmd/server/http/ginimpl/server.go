package ginimpl

import (
	"github.com/995933447/eventobserver"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	. "github.com/vision-first/wegod/internal/app/ginimpl/providers"
	"github.com/vision-first/wegod/internal/pkg/boot"
	. "github.com/vision-first/wegod/internal/pkg/boot/providers"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/event"
	"github.com/vision-first/wegod/internal/pkg/event/handlefuncs"
)

func RunServer() error {
	srv := gin.Default()

	if err := bootServiceProviders(srv); err != nil {
		return err
	}

	if err := srv.Run("0.0.0.0:8081"); err != nil {
		return err
	}

	return nil
}

func bootServiceProviders(srv *gin.Engine) error {
	logger := facades.MustLogger()

	bootstrapper := boot.NewBootstrapper([]boot.ServiceProvider{
		// 数据模型迁移服务提供者
		NewMigrateDataModelProvider([]interface{}{
			&models.Buddha{},
			&models.BuddhaFollow{},
			&models.BuddhaWorship{},
			&models.WorshipProp{},
			&models.User{},
			&models.UserPray{},
			&models.UserDonationDailyStat{},
			&models.DonationOrder{},
			&models.PrayProp{},
			&models.ShopProductCategory{},
			&models.ShopProduct{},
			&models.ShopOrder{},
			&models.Music{},
			&models.PostCategory{},
			&models.Post{},
		}, logger),
		// 事件服务提供者
		NewEventProvider(map[string][]*eventobserver.Listener{
			event.EventNameCreatedDonationOrder: {
				eventobserver.NewListener(event.ListenerNameUserDonationDailyStat, handlefuncs.MakeUserDonationDailyStatHandleFunc(logger)),
			},
		}, logger),
		// http路由服务提供者
		NewHttpRouterProvider(srv),
	})
	if err := bootstrapper.Boot(); err != nil {
		return err
	}
	return nil
}