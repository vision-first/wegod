package services

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"context"
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
	return err
}

func (p *Post) PagePosts(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.Post, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, p.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	// TODO. set option handler
	var postDOs []*models.Post
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &postDOs)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, nil, p.TransErr(err)
	}

	return postDOs, pagination, nil
}