package services

import (
	"github.com/995933447/log-go"
)

type ShopProduct struct {
	logger *log.Logger
}

func NewShopProduct(logger *log.Logger) *ShopProduct {
	return &ShopProduct {
		logger: logger,
	}
}

func (s *ShopProduct) TransErr(err error) error {
	return err
}


