package dtos

import (
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
)

type PagePostCategoriesResp struct {
	PageQueryResp
	// 数据传输对象的作用是解藕响应层和数据层的以来，但是文章分类字段一般不会有什么变更，直接引用数据对象即可，无须过度设计
	List []*models.PostCategory `json:"categories"`
}

type PagePostsResp struct {
	PageQueryResp
	// 数据传输对象的作用是解藕响应层和数据层的以来，但是文章字段一般不会有什么变更，直接引用数据对象即可，无须过度设计
	List []*models.Post `json:"post"`
}


type GetPostReq struct {
	Id uint64
}

type GetPostResp struct {
	Post *models.Post
}