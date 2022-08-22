package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type UserDonationStat struct {
	logger *log.Logger
}

func NewUserDonationStat(logger *log.Logger) *UserDonationStat {
	return &UserDonationStat {
		logger: logger,
	}
}

func (u *UserDonationStat) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}


func (u *UserDonationStat) PageUserDonationRanks(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.UserDonationDailyStat, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, u.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)

	var selectColumns []string
	queryStreamProcessor.
		OnStringList(queryoptions.SelectColumns, func(val []string) error {
			selectColumns = val
			return nil
		}).
		OnTimestampRange(queryoptions.TimeRangeAt, func(beginAt, endAt int64) error {
			db.Where(
				enum.FieldDateYmd + " >= ? AND " + enum.FieldDateYmd + " <= ?",
				time.Unix(beginAt, 0).Format("20060102"),
				time.Unix(endAt, 0).Format("20060102"),
				)
			return nil
		}).OnNone(queryoptions.OrderByPayedMoneyDesc, func() error {
			db.Order(clause.OrderByColumn{
				Column: clause.Column{
					Name: enum.FieldTotalPayedMoney,
				},
				Desc: true,
			})
			return nil
		}).OnNone(queryoptions.GroupByUserId, func() error {
			db.Group(enum.FieldUserId)

			for i, column := range selectColumns {
				if column != enum.FieldTotalPayedMoney {
					continue
				}

				nextIdx := i + 1
				if nextIdx >= len(selectColumns) {
					selectColumns = selectColumns[:i]
					break
				}

				selectColumns = append(selectColumns[:i], selectColumns[nextIdx:]...)
				break
			}
			selectColumns = append(selectColumns, "SUM(" + enum.FieldTotalPayedMoney + ") AS " + enum.FieldTotalPayedMoney)
			return nil
		})

	if len(selectColumns) > 0 {
		db.Select(strings.Join(selectColumns, ","))
	}

	var userDonationDailyStatDOs []*models.UserDonationDailyStat
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &userDonationDailyStatDOs)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, nil, u.TransErr(err)
	}

	return userDonationDailyStatDOs, pagination, nil
}