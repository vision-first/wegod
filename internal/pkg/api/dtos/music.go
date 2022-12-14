package dtos

type Music struct {
	AudioUrl string `json:"audio_url"`
	Name string `json:"name"`
}

type PageMusicsResp struct {
	PageQueryResp
	List []*Music
}