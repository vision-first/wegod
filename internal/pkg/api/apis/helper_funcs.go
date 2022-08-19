package apis

import (
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
)

func transBuddhaDOToDTO(buddhaDO *models.Buddha) *dtos.Buddha {
	return &dtos.Buddha{
		Id: buddhaDO.Id,
		Name: buddhaDO.Name,
		Image: buddhaDO.Image,
		Sort: buddhaDO.Sort,
	}
}

func transBuddhaDOsToDTOs(buddhaDOs []*models.Buddha) []*dtos.Buddha {
	var buddhaDTOs []*dtos.Buddha
	for _, buddhaDO := range buddhaDOs {
		buddhaDTOs = append(buddhaDTOs, transBuddhaDOToDTO(buddhaDO))
	}
	return buddhaDTOs
}

func transBuddhaRentPackageDOToDTO(buddhaRentPackageDo *models.BuddhaRentPackage) *dtos.BuddhaRentPackage {
	return &dtos.BuddhaRentPackage{
		Id: buddhaRentPackageDo.Id,
		Name: buddhaRentPackageDo.Name,
		Desc: buddhaRentPackageDo.Desc,
		Price: buddhaRentPackageDo.Price,
		AvailableDuration: buddhaRentPackageDo.AvailableDuration,
		ShelfStatus: buddhaRentPackageDo.ShelfStatus,
		BuddhaId: buddhaRentPackageDo.BuddhaId,
	}
}

func transBuddhaRentPackageDOsToDTOs(buddhaRentPackageDOs []*models.BuddhaRentPackage) []*dtos.BuddhaRentPackage {
	var buddhaRentPackageDTOs []*dtos.BuddhaRentPackage
	for _, buddhaRentPackageDO := range buddhaRentPackageDOs {
		buddhaRentPackageDTOs = append(buddhaRentPackageDTOs, transBuddhaRentPackageDOToDTO(buddhaRentPackageDO))
	}
	return buddhaRentPackageDTOs
}

func transBuddhaRentOrderDOToDTO(buddhaRentOrderDO *models.BuddhaRentOrder) *dtos.BuddhaRentOrder {
	return &dtos.BuddhaRentOrder{
		OrderSn: buddhaRentOrderDO.Sn,
		BuddhaId: buddhaRentOrderDO.BuddhaId,
		RentPackageId: buddhaRentOrderDO.RentPackageId,
		Price: buddhaRentOrderDO.RentPackagePriceSnapshot,
		RentPackageName: buddhaRentOrderDO.RentPackageNameSnapshot,
		RentPackageDesc: buddhaRentOrderDO.RentPackageDescSnapshot,
		Status: buddhaRentOrderDO.Status,
	}
}

func transPrayPropDOToDTO(prayPropDO *models.PrayProp) *dtos.PrayProp {
	return &dtos.PrayProp{
		Name: prayPropDO.Name,
		Price: prayPropDO.Price,
		Remark: prayPropDO.Remark,
		ShelfStatus: prayPropDO.ShelfStatus,
		Extra: string(prayPropDO.Extra),
	}
}

func transPrayPropDOsToDTOs(prayPropDOs []*models.PrayProp) []*dtos.PrayProp {
	var prayPropDTOs []*dtos.PrayProp
	for _, prayPropDO := range prayPropDOs {
		prayPropDTOs = append(prayPropDTOs, transPrayPropDOToDTO(prayPropDO))
	}
	return prayPropDTOs
}

func transWishPropDOsToDTOs(worshipPropDOs []*models.WorshipProp) []*dtos.WorshipProp {
	var worshipPropDTOs []*dtos.WorshipProp
	for _, worshipPropDO := range worshipPropDOs {
		worshipPropDTOs = append(worshipPropDTOs, &dtos.WorshipProp{
			Id: worshipPropDO.Id,
			Name: worshipPropDO.Name,
			Image: worshipPropDO.Image,
			Price: worshipPropDO.Price,
			ShelfStatus: worshipPropDO.ShelfStatus,
			AvailableDuration: worshipPropDO.AvailableDuration,
		})
	}
	return worshipPropDTOs
}

func transShopProductDOToDTO(shopProductDO *models.ShopProduct) *dtos.ShopProduct {
	return &dtos.ShopProduct{
		Name: shopProductDO.Name,
		CategoryId: shopProductDO.CategoryId,
		MainImages: shopProductDO.MainImages,
		Desc: shopProductDO.Desc,
		ProductType: shopProductDO.ProductType,
		DeliveryType: shopProductDO.DeliveryType,
		ShelfStatus: shopProductDO.ShelfStatus,
		InventoryNum: shopProductDO.InventoryNum,
		SaleNum: shopProductDO.SaleNum,
		Price: shopProductDO.Price,
	}
}

func transShopProductDOsToDTOs(shopProductDOs []*models.ShopProduct) []*dtos.ShopProduct {
	var shopProductDTOs []*dtos.ShopProduct
	for _, shopProductDO := range shopProductDOs {
		shopProductDTOs = append(shopProductDTOs, transShopProductDOToDTO(shopProductDO))
	}
	return shopProductDTOs
}