package providers

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

type MigrateDataModelProvider struct {
	dataModels []interface{}
	logger *log.Logger
	ctx context.Context
}

func NewMigrateDataModelProvider(dataModels []interface{}, logger *log.Logger) *MigrateDataModelProvider {
	return &MigrateDataModelProvider{
		dataModels: dataModels,
		logger: logger,
	}
}

func (p *MigrateDataModelProvider) Boot() error {
	if err := facades.MustGormDB(nil, p.logger).AutoMigrate(p.dataModels...); err != nil {
		return err
	}
	return nil
}
