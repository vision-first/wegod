package servicegen

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

func GenSrv(srvName, servicesPkgCodePath string) error {
	alreadyExistSrvStruct, codeFilePath, err := checkExistSrvStructAndGetCodeFilePath(srvName, servicesPkgCodePath)
	if err != nil {
		return err
	}

	if alreadyExistSrvStruct {
		return fmt.Errorf("service struct:%s already exist", srvName)
	}

	var isCodeFilePathEmpty bool
	if _, err = os.Stat(codeFilePath); os.IsNotExist(err) {
		isCodeFilePathEmpty = true
	}

	apiFp, err := os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return err
	}

	if !isCodeFilePathEmpty {
		fileContent, err := ioutil.ReadAll(apiFp)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return err
		}
		if !file.Package.IsValid() {
			isCodeFilePathEmpty = true
		}
	}

	outputTemplate := "\n\n"

	if !alreadyExistSrvStruct {
		if isCodeFilePathEmpty {
			outputTemplate = `package services

import (
	"github.com/995933447/log-go"
)

`
		}

		tmpl := template.New("srvStruct")
		tmpl, err = tmpl.Parse(srvTmpl)
		if err != nil {
			return err
		}

		buf := bytes.NewBufferString(outputTemplate)
		err = tmpl.Execute(buf, &struct {
			SrvName string
			SrvNameAbbreviation string
		}{
			SrvName: srvName,
			SrvNameAbbreviation: strings.ToLower(string(srvName[0])),
		})
		if err != nil {
			return err
		}

		outputTemplate = buf.String() + "\n\n"
	}

	outputTemplateBytes := len(outputTemplate)
	var writtenOutputTemplateBytes int
	for {
		n, err := apiFp.Write([]byte(outputTemplate))
		if err != nil {
			return err
		}
		writtenOutputTemplateBytes += n
		if writtenOutputTemplateBytes >= outputTemplateBytes {
			break
		}
		outputTemplate = outputTemplate[writtenOutputTemplateBytes:]
	}

	return nil
}

func checkExistSrvStructAndGetCodeFilePath(srvName, servicesPkgCodePath string) (existSrvStruct bool, codeFilePath string, err error) {
	srvImplsPkgMap, err := parser.ParseDir(token.NewFileSet(), servicesPkgCodePath, nil, 0)
	if err != nil {
		return
	}

	for _, pkg:= range srvImplsPkgMap {
		for filePath, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := typeSpec.Type.(*ast.StructType); ok && typeSpec.Name.String() == srvName {
								existSrvStruct = true
								codeFilePath = filePath
							}
						}
					}
				}
			}
		}
	}

	if codeFilePath == "" {
		codeFilePath = servicesPkgCodePath + "/" + stringhelper.Snake(srvName) + ".go"
	}

	return
}
