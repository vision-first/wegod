package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Shop struct {
	logger *log.Logger
}

func (s *Shop) PageProducts(ctx api.Context, req *dtos.PageShopProductsReq) (*dtos.PageShopProductsResp, error) {
	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Limit)
	queryStream.SetOption(queryoptions.OnShelfStatus, nil)

	var (
		resp dtos.PageShopProductsResp
		shopProductDOs []*models.ShopProduct
		err error
	)
	shopProductDOs, resp.Pagination, err = services.NewShopProduct(s.logger).PageProducts(ctx, queryStream)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transShopProductDOsToDTOs(shopProductDOs)

    return &resp, nil
}

func (s *Shop) CreateOrder(ctx api.Context, req *dtos.CreateShopOrderReq) (*dtos.CreateShopOrderResp, error) {
    var resp dtos.CreateShopOrderResp

    // TODO.write your logic

    return &resp, nil
}

func (s *Shop) PageOrders(ctx api.Context, req *dtos.PageShopOrdersReq) (*dtos.PageShopOrdersResp, error) {
    var resp dtos.PageShopOrdersResp

    // TODO.write your logic

    return &resp, nil
}

func (s *Shop) PageProductCategories(ctx api.Context, req *dtos.PageShopProductCategoriesReq) (*dtos.PageShopProductCategoriesResp, error) {
    var (
    	resp dtos.PageShopProductCategoriesResp
    	productCategoryDOs []*models.ShopProductCategory
		err error
	)
	productCategoryDOs, resp.Pagination, err = services.NewShopProductCategory(s.logger).
		PageProductCategories(ctx, optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset))
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}
	for _, productCategoryDO := range productCategoryDOs {
		resp.List = append(resp.List, &dtos.ShopProductCategory{
			Name: productCategoryDO.Name,
			Id: productCategoryDO.Id,
		})
	}
    return &resp, nil
}