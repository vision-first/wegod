package dtos

import "github.com/995933447/optionstream"

type Buddha struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Sort uint32 `json:"sort"`
}

type PageBuddhaReq struct {
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type PageBuddhaResp struct {
	List []*Buddha `json:"list"`
	Pagination *optionstream.Pagination `json:"pagination"`
}

type WatchBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
}

type WatchBuddhaResp struct {
}

type UnwatchBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
}

type UnwatchBuddhaResp struct {
}

type PageUserWatchedBuddhaReq struct {
	QueryOptions []*optionstream.Option `json:"query_options"`
	Limit        int64                  `json:"limit"`
	Offset       int64                  `json:"offset"`
}

type PageUserWatchedBuddhaResp struct {
	List []*Buddha
	Pagination *optionstream.Pagination
}