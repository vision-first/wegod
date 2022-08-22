package apigen

const (
	apiMethodTmpl = `func ({{.ApiModuleAbbreviation}} *{{.ApiModule}}) {{.ApiMethod}}(ctx api.Context, req *dtos.{{.ApiMethod}}Req) (*dtos.{{.ApiMethod}}Resp, error) {
    var resp dtos.{{.ApiMethod}}Resp

    // TODO.write your logic

    return &resp, nil
}`
	apiStructTmpl = `type {{.ApiModule}} struct {
	logger *log.Logger
}

func New{{.ApiModule}}(logger *log.Logger) *{{.ApiModule}} {
	return &{{.ApiModule}} {
		logger: logger,
	}
}
`
	apiReqDTOTmpl = `type {{.ApiMethod}}Req struct {
}`
	apiRespDTOTmpl = `type {{.ApiMethod}}Resp struct {
}`
)
