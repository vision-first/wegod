package services

import (
	"github.com/vision-first/wegod/internal/pkg/datamodels"
)

type Buddha interface {
	pageBuddhas() ([]*datamodels.Buddha, error)
}