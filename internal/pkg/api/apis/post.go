package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
)

type Post struct {
	logger log.Logger
}

func (p *Post) PageCategories(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageCategoriesResp, error) {
    var resp dtos.PageCategoriesResp

    // TODO.Write your logic

    return &resp, nil
}

func (p *Post) PagePosts(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PagePostsResp, error) {
    var resp dtos.PagePostsResp

    // TODO.Write your logic

    return &resp, nil
}