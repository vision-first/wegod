package ginimpl

import (
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/facades"
	. "github.com/vision-first/wegod/internal/app/ginimpl/providers"
	"github.com/vision-first/wegod/internal/pkg/boot"
	. "github.com/vision-first/wegod/internal/pkg/boot/providers"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
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
			&datamodels.Buddha{},
			&datamodels.BuddhaFollow{},
			&datamodels.BuddhaWorship{},
			&datamodels.BuddhaWorshipProp{},
			&datamodels.User{},
			&datamodels.UserPray{},
			&datamodels.UserDonationDailyStat{},
			&datamodels.UserDonationRecord{},
			&datamodels.PrayProp{},
			&datamodels.ShopProductCategory{},
			&datamodels.ShopProduct{},
			&datamodels.ShopOrder{},
			&datamodels.Music{},
			&datamodels.PostCategory{},
			&datamodels.Post{},
		}, facades.MustLogger()),
		NewHttpRouterProvider(srv),
	})
	if err := bootstrapper.Boot(); err != nil {
		return err
	}
	return nil
}