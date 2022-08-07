package apigen

import (
	"bytes"
	"fmt"
	"github.com/995933447/stringhelper-go"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func GenApi(apiModule, apiMethod, apiPkgPath, templatesDirPath string) error {
	apiPkgPath = strings.TrimRight(apiPkgPath, "/")
	apiImplsPkgPath := apiPkgPath + "/apis"
	dtosPkgPath := apiPkgPath + "/dtos"
	templatesDirPath = strings.TrimRight(templatesDirPath, "/")

	fset := token.NewFileSet()
	apiImplsPkgMap, err := parser.ParseDir(fset, apiImplsPkgPath, nil, 0)
	if err != nil {
		return err
	}

	var (
		alreadyGenApiStruct bool
		apiCodeFilePath string
	)
	for _, pkg:= range apiImplsPkgMap {
		for filePath, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						//if importSpec, ok := spec.(*ast.ImportSpec); ok {
						//	if importSpec.Path == nil {
						//		continue
						//	}
						//
						//	if importSpec.Path.Value == apiDTOsPkg {
						//		alreadyImportedDTOsPkgInFile = true
						//	}
						//}

						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := typeSpec.Type.(*ast.StructType); ok && typeSpec.Name.String() == apiModule {
								alreadyGenApiStruct = true
								apiCodeFilePath = filePath
							}
						}
					}
				}

				if funcDecl, ok := decl.(*ast.FuncDecl); ok &&
					funcDecl.Name.String() == apiMethod &&
					funcDecl.Recv != nil &&
					len(funcDecl.Recv.List) == 1 {
					for _, receiverIdent := range funcDecl.Recv.List[0].Names {
						if receiverIdent.String() == apiModule {
							return fmt.Errorf("api method: %s.%s already exist", apiModule, apiMethod)
						}
					}
				}
			}
		}
	}

	dtosPkgMap, err := parser.ParseDir(fset, dtosPkgPath, nil,0)
	var (
		alreadyExistReqDTO bool
		alreadyExistRespDTO bool
		dtoCodeFilePath string
	)
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
							alreadyExistReqDTO = true
							dtoCodeFilePath = filePath
						case apiMethod + "Resp":
							alreadyExistReqDTO = true
							dtoCodeFilePath = filePath
						}
					}
				}
			}
		}
	}

	var isDTOCodeFileEmpty bool
	if dtoCodeFilePath == "" {
		dtoCodeFilePath = apiPkgPath + "/dtos/" + apiModule + ".go"
		isDTOCodeFileEmpty = true
	} else {
		file, err := parser.ParseFile(fset, dtoCodeFilePath, "", 0)
		if err != nil {
			return err
		}
		if !file.Package.IsValid() {
			isDTOCodeFileEmpty = true
		}
	}

	var outputDTOsTemplate string
	if !alreadyExistReqDTO {
		tmpl := template.New("apiReqDTO")
		tmpl, err := tmpl.ParseFiles(templatesDirPath + "/api_req_dto.tmpl")
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
		outputDTOsTemplate = buf.String() + "\n"
	}

	if !alreadyExistRespDTO {
		tmpl := template.New("apiRespDTO")
		tmpl, err := tmpl.ParseFiles(templatesDirPath + "/api_resp_dto.tmpl")
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
		outputDTOsTemplate = buf.String() + "\n"
	}

	var dtosFp *os.File
	if isDTOCodeFileEmpty {
		outputDTOsTemplate = `package dtos
` + outputDTOsTemplate
		dtosFp, err = os.OpenFile(dtoCodeFilePath, os.O_CREATE|os.O_RDWR, 7555)
		if err != nil {
			return err
		}
	} else {
		dtosFp, err = os.OpenFile(dtoCodeFilePath, os.O_RDWR, 755)
		if err != nil {
			return err
		}
		fileContent, err := ioutil.ReadAll(dtosFp)
		if err != nil {
			return err
		}
		outputDTOsTemplate = string(fileContent) + "\n" + outputDTOsTemplate
	}

	outputDTOsTemplateBytes := len(outputDTOsTemplate)
	var writtenOutputDTOsTemplateBytes int
	for {
		n, err := dtosFp.Write([]byte(outputDTOsTemplate))
		if err != nil {
			return err
		}
		writtenOutputDTOsTemplateBytes += n
		if writtenOutputDTOsTemplateBytes >= outputDTOsTemplateBytes {
			break
		}
	}

	if apiCodeFilePath == "" {
		apiCodeFilePath = apiImplsPkgPath + "/" + stringhelper.Snake(apiModule) + ".go"
	}

	var (
		outputApiTemplate string
		apiFp *os.File
	)
	if alreadyGenApiStruct {
		if apiFp, err = os.OpenFile(apiCodeFilePath, os.O_RDWR, 0755); err != nil {
			return err
		}

		fileContent, err := ioutil.ReadAll(apiFp)
		if err != nil {
			return err
		}

		outputApiTemplate += string(fileContent) + "\n"
	} else {
		outputApiTemplate = `package apis

`
		tmpl := template.New("apiStruct")
		tmpl, err = tmpl.ParseFiles(templatesDirPath + "/api_struct.tmpl")
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

		if apiFp, err = os.OpenFile(apiCodeFilePath, os.O_CREATE|os.O_RDWR, 755); err != nil {
			return err
		}
	}

	tmpl := template.New("apiMethod")
	tmpl, err = tmpl.ParseFiles(templatesDirPath + "/api_method.tmpl")
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

	return nil
}
