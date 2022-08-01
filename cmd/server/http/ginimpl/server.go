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

	bootstrapper := boot.NewBootstrapper([]boot.ServiceProvider{
		NewMigrateDataModelProvider([]interface{}{
			&datamodels.Buddha{},
		}, facades.MustLogger()),
		NewHttpRouterProvider(srv),
	})
	if err := bootstrapper.Boot(); err != nil {
		return err
	}

	if err := srv.Run("0.0.0.0:8081"); err != nil {
		return err
	}

	return nil
}