package services

import (
	"github.com/995933447/log-go"
	"gorm.io/gorm"
)

type ShopOrder struct {
	logger *log.Logger
}

func NewShopOrder(logger *log.Logger) *ShopOrder {
	return &ShopOrder {
		logger: logger,
	}
}

func (s *ShopOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}


