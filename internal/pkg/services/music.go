package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
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
	db := facades.MustGormDB(ctx, m.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	var musicDOs []*models.Music
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &musicDOs)
	if err != nil {
		m.logger.Error(ctx, err)
		return nil, nil, m.TransErr(err)
	}

	return musicDOs, pagination, nil
}