package providers

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/boot"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

type MigrateDataModelProvider struct {
	dataModels []interface{}
	logger *log.Logger
	ctx context.Context
}

var _ boot.ServiceProvider = (*MigrateDataModelProvider)(nil)

func NewMigrateDataModelProvider(dataModels []interface{}, logger *log.Logger) *MigrateDataModelProvider {
	return &MigrateDataModelProvider{
		dataModels: dataModels,
		logger: logger,
	}
}

func (p *MigrateDataModelProvider) Boot() error {
	if err := facades.MustGORMDB(nil, p.logger).AutoMigrate(p.dataModels...); err != nil {
		return err
	}
	return nil
}
