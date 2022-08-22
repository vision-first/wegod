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

func (s *Shop) PageProducts(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageShopProductsResp, error) {
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

    authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

	orderDO, err := services.NewShopOrder(s.logger).CreateOrder(
		ctx,
		authIdent.GetUserId(),
		req.ProductId,
		req.Num,
		req.Remark,
		services.Consignee{
			ConsigneeName: req.ConsigneeName,
			ConsigneePhone: req.ConsigneePhone,
			ConsigneeProvince: req.ConsigneeProvince,
			ConsigneeCity: req.ConsigneeCity,
			ConsigneeDistrict: req.ConsigneeDistrict,
			ConsigneeAddress: req.ConsigneeAddress,
			LogisticsSn: req.LogisticsSn,
			LogisticsCompany: req.LogisticsCompany,
		},
		)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

	resp.Order = transShopOrderDOToDTO(orderDO)

    return &resp, nil
}

func (s *Shop) PageOrders(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageShopOrdersResp, error) {
    var resp dtos.PageShopOrdersResp

	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)

	var (
		orderDOs []*models.ShopOrder
		err error
	)
	orderDOs, resp.Pagination, err = services.NewShopOrder(s.logger).PageOrders(ctx, queryStream)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transShopOrderDOsToDTOs(orderDOs)

    return &resp, nil
}

func (s *Shop) PageProductCategories(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageShopProductCategoriesResp, error) {
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

func (s *Shop) GetShopOrder(ctx api.Context, req *dtos.GetOrderReq) (*dtos.GetOrderResp, error) {
    var resp dtos.GetOrderResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

    orderDO, err := services.NewShopOrder(s.logger).GetOrder(
		ctx,
		optionstream.NewStream(nil).
			SetOption(queryoptions.EqualId, req.OrderId).
			SetOption(queryoptions.EqualUserId, authIdent.GetUserId()),
		)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, err
	}

	resp.Order = transShopOrderDOToDTO(orderDO)

	return &resp, nil
}