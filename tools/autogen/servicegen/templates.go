package servicegen

const (
	srvTmpl = `
type {{.SrvName}} struct {
	logger *log.Logger
}

func ({{.SrvNameAbbreviation}} *{{.SrvName}}) TransErr(err error) error {
	switch err {
	
	}
	return err
}
`
)
