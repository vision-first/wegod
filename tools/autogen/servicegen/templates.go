package servicegen

const (
	srvTmpl = `type {{.Service}} struct {
	logger *log.Logger
}

func New{{.Service}}(logger *log.Logger) *{{.Service}} {
	return &{{.Service}} {
		logger: logger,
	}
}

func ({{.ServiceAbbreviation}} *{{.Service}}) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}
`

	optionStreamPageQueryTmpl = `func ({{.ServiceAbbreviation}} *{{.Service}}) {{.Method}}(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.{{.Model}}, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, {{.ServiceAbbreviation}}.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	// TODO. set option handler
	var {{.LowerFirstCharModel}}DOs []*models.{{.Model}}
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &{{.LowerFirstCharModel}}DOs)
	if err != nil {
		{{.ServiceAbbreviation}}.logger.Error(ctx, err)
		return nil, nil, {{.ServiceAbbreviation}}.TransErr(err)
	}

	return {{.LowerFirstCharModel}}DOs, pagination, nil
}`
)
