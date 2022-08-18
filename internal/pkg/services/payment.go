package services

import (
	"github.com/995933447/log-go"
	"gorm.io/gorm"
)

type Payment struct {
	logger *log.Logger
}

func NewPayment(logger *log.Logger) *Payment {
	return &Payment {
		logger: logger,
	}
}

func (p *Payment) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}


