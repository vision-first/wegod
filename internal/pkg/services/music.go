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
)

type Music struct {
	logger *log.Logger
}

func NewMusic(logger *log.Logger) *Music {
	return &Music {
		logger: logger,
	}
}

func (m *Music) TransErr(err error) error {
	return err
}
func (m *Music) PageMusics(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.Music, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, m.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnNone(queryoptions.OnShelfStatus, func() error {
			db.Where(&models.Music{ShelfStatus: models.ShelfStatusOnShelf})
			return nil
		}).
		OnUint64(queryoptions.EqualBuddhaId, func(val uint64) error {
			db.Where("JSON_CONTAINS(" + enum.FieldBuddhaIds + ", ?)", val)
			return nil
		})

	var musicDOs []*models.Music
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &musicDOs)
	if err != nil {
		m.logger.Error(ctx, err)
		return nil, nil, m.TransErr(err)
	}

	return musicDOs, pagination, nil
}