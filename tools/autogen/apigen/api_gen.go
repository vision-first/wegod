package apigen

import (
	"bytes"
	"fmt"
	"github.com/995933447/stringhelper-go"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func GenApi(apiModule, apiMethod, apiPkgPath, apiPkgCodePath string) error {
	if err := genApiDTOs(apiModule, apiMethod, apiPkgCodePath); err != nil {
		return err
	}

	if err := genApiModuleAndMethod(apiModule, apiMethod, apiPkgPath, apiPkgCodePath); err != nil {
		return err
	}

	return nil
}

func genApiModuleAndMethod(apiModule, apiMethod, apiPkgPath, apiPkgCodePath string) error {
	alreadyExistApiStruct, alreadyExistApiMethod, codeFilePath, err := checkExistApiOrMethodAndGetCodePath(apiModule, apiMethod, apiPkgCodePath)
	if err != nil {
		return err
	}

	if alreadyExistApiMethod {
		return fmt.Errorf("api method: %s.%s already exist", apiModule, apiMethod)
	}

	var isApiCodeFilePathEmpty bool
	if _, err = os.Stat(codeFilePath); os.IsNotExist(err) {
		isApiCodeFilePathEmpty = true
	}

	apiFp, err := os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return err
	}

	if !isApiCodeFilePathEmpty {
		fileContent, err := ioutil.ReadAll(apiFp)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return err
		}
		if !file.Package.IsValid() {
			isApiCodeFilePathEmpty = true
		}
	}

	outputApiTemplate := "\n\n"

	if !alreadyExistApiStruct {
		if isApiCodeFilePathEmpty {
			outputApiTemplate = `package apis

import (
	"github.com/995933447/log-go"
	"` + apiPkgPath + `"
	"` + apiPkgPath + `/dtos"
)

`
		}

		tmpl := template.New("apiStruct")
		tmpl, err = tmpl.Parse(apiStructTmpl)
		if err != nil {
			return err
		}

		buf := bytes.NewBufferString(outputApiTemplate)
		err = tmpl.Execute(buf, &struct {
			ApiModule string
		}{
			ApiModule: apiModule,
		})
		if err != nil {
			return err
		}

		outputApiTemplate = buf.String() + "\n\n"
	}

	tmpl := template.New("apiMethod")
	tmpl, err = tmpl.Parse(apiMethodTmpl)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(outputApiTemplate)
	err = tmpl.Execute(buf, &struct {
		ApiModuleAbbreviation string
		ApiModule string
		ApiMethod string
	}{
		ApiModuleAbbreviation: strings.ToLower(string(apiModule[0])),
		ApiModule: apiModule,
		ApiMethod: apiMethod,
	})
	if err != nil {
		return err
	}
	outputApiTemplate = buf.String()

	outputApiTemplateBytes := len(outputApiTemplate)
	var writtenOutputApiTemplateBytes int
	for {
		n, err := apiFp.Write([]byte(outputApiTemplate))
		if err != nil {
			return err
		}
		writtenOutputApiTemplateBytes += n
		if writtenOutputApiTemplateBytes >= outputApiTemplateBytes {
			break
		}
		outputApiTemplate = outputApiTemplate[writtenOutputApiTemplateBytes:]
	}

	return nil
}

func checkExistApiOrMethodAndGetCodePath(apiModule, apiMethod, apiPkgCodePath string) (existApiStruct, existApiMethod bool, codeFilePath string, err error) {
	apiPkgCodePath = strings.TrimRight(apiPkgCodePath, "/")
	apiImplsPkgCodePath := apiPkgCodePath + "/apis"

	apiImplsPkgMap, err := parser.ParseDir(token.NewFileSet(), apiImplsPkgCodePath, nil, 0)
	if err != nil {
		return
	}

	for _, pkg:= range apiImplsPkgMap {
		for filePath, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := typeSpec.Type.(*ast.StructType); ok && typeSpec.Name.String() == apiModule {
								existApiStruct = true
								codeFilePath = filePath
							}
						}
					}
				}

				if funcDecl, ok := decl.(*ast.FuncDecl); ok &&
					funcDecl.Name.String() == apiMethod &&
					funcDecl.Recv != nil &&
					len(funcDecl.Recv.List) == 1 {
					starExprRecvType, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
					if !ok {
						continue
					}
					starExprRecvTypeIdent, ok := starExprRecvType.X.(*ast.Ident)
					if !ok {
						continue
					}
					if starExprRecvTypeIdent.String() == apiModule {
						existApiMethod = true
					}
				}
			}
		}
	}

	if codeFilePath == "" {
		codeFilePath = apiImplsPkgCodePath + "/" + stringhelper.Snake(apiModule) + ".go"
	}

	return
}

func genApiDTOs(apiModule, apiMethod, apiPkgCodePath string) error {
	alreadyExistReqDTO, alreadyExistRespDTO, dtoCodeFilePath, err := checkExistDTOsAndGetCodeFilePath(apiModule, apiMethod, apiPkgCodePath)
	if err != nil {
		return err
	}

	var isDTOCodeFileEmpty bool
	if _, err = os.Stat(dtoCodeFilePath); os.IsNotExist(err) {
		isDTOCodeFileEmpty = true
	}

	dtoCodeFp, err := os.OpenFile(dtoCodeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return err
	}

	if !isDTOCodeFileEmpty {
		fileContent, err := ioutil.ReadAll(dtoCodeFp)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return err
		}
		if !file.Package.IsValid() {
			isDTOCodeFileEmpty = true
		}
	}

	outputDTOsTemplate := "\n\n"
	if isDTOCodeFileEmpty {
		outputDTOsTemplate = `package dtos

`
	}

	if !alreadyExistReqDTO {
		tmpl := template.New("apiReqDTO")
		tmpl, err := tmpl.Parse(apiReqDTOTmpl)
		if err != nil {
			return err
		}
		buf := bytes.NewBufferString(outputDTOsTemplate)
		err = tmpl.Execute(buf, &struct {
			ApiMethod string
		}{
			ApiMethod: apiMethod,
		})
		if err != nil {
			return err
		}
		outputDTOsTemplate = buf.String() + "\n\n"
	}

	if !alreadyExistRespDTO {
		tmpl := template.New("apiRespDTO")
		tmpl, err := tmpl.Parse(apiRespDTOTmpl)
		if err != nil {
			return err
		}
		buf := bytes.NewBufferString(outputDTOsTemplate)
		err = tmpl.Execute(buf, &struct {
			ApiMethod string
		}{
			ApiMethod: apiMethod,
		})
		if err != nil {
			return err
		}
		outputDTOsTemplate = buf.String()
	}

	outputDTOsTemplateBytes := len(outputDTOsTemplate)
	var writtenOutputDTOsTemplateBytes int
	for {
		n, err := dtoCodeFp.Write([]byte(outputDTOsTemplate))
		if err != nil {
			return err
		}
		writtenOutputDTOsTemplateBytes += n
		if writtenOutputDTOsTemplateBytes >= outputDTOsTemplateBytes {
			break
		}
		outputDTOsTemplate = outputDTOsTemplate[writtenOutputDTOsTemplateBytes:]
	}

	return nil
}

func checkExistDTOsAndGetCodeFilePath(apiModule, apiMethod, apiPkgCodePath string) (existReqDTO, existRespDTO bool, codeFilePath string, err error) {
	dtosCodePkgPath := apiPkgCodePath + "/dtos"

	dtosPkgMap, err := parser.ParseDir(token.NewFileSet(), dtosCodePkgPath, nil,0)
	if err != nil {
		return
	}

	for _, pkg := range dtosPkgMap {
		for filePath, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				if genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						structTypeName := typeSpec.Name.String()
						switch structTypeName {
						case apiMethod + "Req":
							existReqDTO = true
							codeFilePath = filePath
						case apiMethod + "Resp":
							existRespDTO = true
							codeFilePath = filePath
						}
					}
				}
			}
		}
	}

	if codeFilePath == "" {
		codeFilePath = apiPkgCodePath + "/dtos/" + stringhelper.Snake(apiModule) + ".go"
	}

	return existReqDTO, existRespDTO, codeFilePath, nil
}