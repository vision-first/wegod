package apigen

const (
	apiMethodTmpl = `func ({{.ApiModuleAbbreviation}} *{{.ApiModule}}) {{.ApiMethod}}(ctx api.Context, req *dtos.{{.ApiMethod}}Req) (*dtos.{{.ApiMethod}}Resp, error) {
    var resp dtos.{{.ApiMethod}}Resp

    // TODO.Write your logic

    return &resp, nil
}`
	apiStructTmpl = `type {{.ApiModule}} struct {
	logger *log.Logger
}`
	apiReqDTOTmpl = `type {{.ApiMethod}}Req struct {
}`
	apiRespDTOTmpl = `type {{.ApiMethod}}Resp struct {
}`
)
