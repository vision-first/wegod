package ginimpl

import (
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	. "github.com/vision-first/wegod/internal/app/ginimpl/providers"
	"github.com/vision-first/wegod/internal/pkg/boot"
	. "github.com/vision-first/wegod/internal/pkg/boot/providers"
	"github.com/vision-first/wegod/internal/pkg/models"
	"github.com/vision-first/wegod/internal/pkg/event"
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
	bootstrapper := boot.NewBootstrapper([]boot.ServiceProvider{
		NewMigrateDataModelProvider([]interface{}{
			&models.Buddha{},
			&models.BuddhaFollow{},
			&models.BuddhaWorship{},
			&models.BuddhaWorshipProp{},
			&models.User{},
			&models.UserPray{},
			&models.UserDonationDailyStat{},
			&models.UserDonationRecord{},
			&models.PrayProp{},
			&models.ShopProductCategory{},
			&models.ShopProduct{},
			&models.ShopOrder{},
			&models.Music{},
			&models.PostCategory{},
			&models.Post{},
		}, facades.MustLogger()),
		NewEventProvider(map[string][]*event.Listener{

		}, facades.MustLogger()),
		NewHttpRouterProvider(srv),
	})
	if err := bootstrapper.Boot(); err != nil {
		return err
	}
	return nil
}