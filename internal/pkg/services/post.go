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
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

type Post struct {
	logger *log.Logger
}

func NewPost(logger *log.Logger) *Post {
	return &Post {
		logger: logger,
	}
}

func (p *Post) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodePostNotFound)
	}
	return err
}

func (p *Post) PagePosts(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.Post, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, p.logger).Order(clause.OrderByColumn{
		Column: clause.Column{Name: enum.FieldSort},
		Desc: true,
	})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.OnUint64(queryoptions.EqualCategoryId, func(val uint64) error {
		db.Where(&models.Post{CategoryId: val})
		return nil
	})

	var postDOs []*models.Post
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &postDOs)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, nil, p.TransErr(err)
	}

	return postDOs, pagination, nil
}

func (p *Post) GetPostById(ctx context.Context, id uint64) (*models.Post, error) {
	var post models.Post
	if err := facades.MustGormDB(ctx, p.logger).First(&post, id).Error; err != nil {
		p.logger.Error(ctx, err)
		return nil, p.TransErr(err)
	}
	return &post, nil
}