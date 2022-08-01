package services

import (
	"wegod/internal/pkg/datamodels"
)

type Buddha interface {
	pageBuddhas() ([]*datamodels.Buddha, error)
}