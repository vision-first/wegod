package handlefuncs

import (
	"context"
	"github.com/995933447/eventobserver"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/event/payloads"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func MakeUserDonationDailyStatHandleFunc(logger *log.Logger) eventobserver.HandleEventFunc {
	return func (ctx context.Context, event eventobserver.Event) error {
		payload := event.GetData().(*payloads.CreatedDonationOrder)

		dateYmd, err := strconv.ParseInt(time.Unix(payload.CreatedAt, 0).Format("20060102"), 0, 32)
		if err != nil {
			logger.Error(ctx, err)
			return nil
		}

		db := facades.MustGormDB(ctx, logger)

		updateStatRes := db.Where(&models.UserDonationDailyStat{UserId: payload.UserId, DateYmd: int(dateYmd)}).
		Updates(map[string]interface{}{
			enum.FieldTotalMoney: gorm.Expr(enum.FieldTotalMoney + " + ?", payload.Money),
			enum.FieldTotalOrderNum: gorm.Expr(enum.FieldTotalOrderNum + " + 1"),
		})
		if updateStatRes.Error != nil {
			logger.Error(ctx, err)
			return nil
		}

		statDO := models.UserDonationDailyStat{
			UserId: payload.UserId,
			TotalMoney: payload.Money,
			TotalOrderNum: 1,
			CalculatedAt: time.Now().Unix(),
			DateYmd: int(dateYmd),
		}
		firstOrCreateRes := db.Where(&models.UserDonationDailyStat{UserId: payload.UserId}).FirstOrCreate(&statDO)
		if firstOrCreateRes.Error != nil {
			logger.Error(ctx, err)
			return nil
		}

		if firstOrCreateRes.RowsAffected > 0 {
			return nil
		}

		err = db.Where(map[string]interface{}{enum.FieldId: statDO.Id}).
		Updates(map[string]interface{}{
			enum.FieldTotalMoney: gorm.Expr(enum.FieldTotalMoney + " + ?", payload.Money),
			enum.FieldTotalOrderNum: gorm.Expr(enum.FieldTotalOrderNum + " + 1"),
		}).
		Error
		if updateStatRes.Error != nil {
			logger.Error(ctx, err)
			return nil
		}

		return nil
	}
}