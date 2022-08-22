package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Post struct {
	logger *log.Logger
}

func NewPost(logger *log.Logger) *Post {
	return &Post{
		logger: logger,
	}
}

func (p *Post) PageCategories(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PagePostCategoriesResp, error) {
	var (
		resp dtos.PagePostCategoriesResp
		err error
	)
	resp.List, resp.Pagination, err = services.NewPostCategory(p.logger).PostCategories(
		ctx,
		optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset),
		)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (p *Post) PagePosts(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PagePostsResp, error) {
    var (
    	resp dtos.PagePostsResp
    	err error
	)
	resp.List, resp.Pagination, err = services.NewPost(p.logger).PagePosts(
		ctx,
		optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset),
	)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (p *Post) GetPost(ctx api.Context, req *dtos.GetPostReq) (*dtos.GetPostResp, error) {
    var (
    	resp dtos.GetPostResp
    	err error
	)
    resp.Post, err = services.NewPost(p.logger).GetPostById(ctx, req.Id)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

	return &resp, nil
}