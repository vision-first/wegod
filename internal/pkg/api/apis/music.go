package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Music struct {
	logger *log.Logger
}

func (m *Music) PageMusics(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageMusicsResp, error) {
    var (
    	resp dtos.PageMusicsResp
    	musicDOs []*models.Music
		err error
	)
    musicDOs, resp.Pagination, err = services.NewMusic(m.logger).PageMusics(ctx, optionstream.NewQueryStream(req.QueryOptions, req.Offset, req.Limit))
	if err != nil {
		m.logger.Error(ctx, err)
		return nil, err
	}

	for _, musicDO := range musicDOs {
		resp.List = append(resp.List, &dtos.Music{
			AudioUrl: musicDO.AudioUrl,
			Name: musicDO.Name,
		})
	}

    return &resp, nil
}