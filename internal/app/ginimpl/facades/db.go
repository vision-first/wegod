package facades

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

func MustGormDB(ctx *gin.Context) *gorm.DB {
	return facades.MustGormDB(MustLogger()).WithContext(ctx)
}
