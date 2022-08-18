package dtos

import "github.com/995933447/optionstream"

type PageQueryReq struct {
	QueryOptions []*optionstream.Option `json:"query_options"`
	Limit        int64                  `json:"limit"`
	Offset       int64                  `json:"offset"`
}

type PageQueryResp struct {
	Pagination *optionstream.Pagination `json:"pagination"`
}
