package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostCategory struct {
	logger *log.Logger
}

func NewPostCategory(logger *log.Logger) *PostCategory {
	return &PostCategory {
		logger: logger,
	}
}

func (p *PostCategory) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeCategoryNotFound)
	}
	return err
}


func (p *PostCategory) PostCategories(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.PostCategory, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, p.logger).Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: enum.FieldSort,
		},
		Desc: true,
	})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	// TODO. set option handler
	var postCategoryDOs []*models.PostCategory
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &postCategoryDOs)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, nil, p.TransErr(err)
	}

	return postCategoryDOs, pagination, nil
}