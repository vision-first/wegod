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
	List []*Buddha
	Pagination *optionstream.Pagination
}

type WatchBuddhaReq struct {
	BuddhaId uint64
}

type WatchBuddhaResp struct {
}

type UnwatchBuddhaReq struct {
}

type UnwatchBuddhaResp struct {
}