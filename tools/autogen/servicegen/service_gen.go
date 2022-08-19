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

func GenService(srvName, servicesPkgCodePath string) error {
	alreadyExistSrvStruct, codeFilePath, err := checkExistSrvStructAndGetCodeFilePath(srvName, servicesPkgCodePath)
	if err != nil {
		return err
	}

	if alreadyExistSrvStruct {
		return fmt.Errorf("service struct:%s already exist", srvName)
	}

	code, err := getSrvStructCodeContent(srvName, codeFilePath)
	if err != nil {
		return err
	}

	codeFp, err := os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return err
	}

	defer codeFp.Close()

	outputTemplateBytes := len(code)
	var writtenOutputTemplateBytes int
	for {
		n, err := codeFp.Write([]byte(code))
		if err != nil {
			return err
		}
		writtenOutputTemplateBytes += n
		if writtenOutputTemplateBytes >= outputTemplateBytes {
			break
		}
		code = code[writtenOutputTemplateBytes:]
	}

	return nil
}

func getSrvStructCodeContent(srvName, codeFilePath string) (string, error) {
	var isCodeFilePathEmpty bool
	if _, err := os.Stat(codeFilePath); os.IsNotExist(err) {
		isCodeFilePathEmpty = true
	}

	if !isCodeFilePathEmpty {
		codeFp, err := os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
		if err != nil {
			return "", err
		}

		defer codeFp.Close()

		fileContent, err := ioutil.ReadAll(codeFp)
		if err != nil {
			return "", err
		}

		file, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return "", err
		}

		if !file.Package.IsValid() {
			isCodeFilePathEmpty = true
		}
	}

	outputTemplate := "\n\n"

	if isCodeFilePathEmpty {
		outputTemplate = `package services

import (
	"github.com/995933447/log-go"
	"gorm.io/gorm"
)

`
	}

	tmpl := template.New("srvStruct")
	tmpl, err := tmpl.Parse(srvTmpl)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString(outputTemplate)
	err = tmpl.Execute(buf, &struct {
		Service string
		ServiceAbbreviation string
	}{
		Service: srvName,
		ServiceAbbreviation: strings.ToLower(string(srvName[0])),
	})
	if err != nil {
		return "", err
	}

	outputTemplate = buf.String() + "\n\n"

	return outputTemplate, nil
}

func getOptionStreamPageQueryMethodCodeContent(srvName, methodName, modelName string) (string, error) {
	tmpl := template.New("srvMethod")
	tmpl, err := tmpl.Parse(optionStreamPageQueryTmpl)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	var lowerFirstCharModel string
	if len(modelName) > 1 {
		lowerFirstCharModel = strings.ToLower(string(modelName[0])) + modelName[1:]
	} else {
		lowerFirstCharModel = strings.ToLower(string(modelName[0]))
	}
	err = tmpl.Execute(buf, &struct {
		Service string
		ServiceAbbreviation string
		Method string
		Model string
		LowerFirstCharModel string
	}{
		Service: srvName,
		ServiceAbbreviation: strings.ToLower(string(srvName[0])),
		Method: methodName,
		Model: modelName,
		LowerFirstCharModel: lowerFirstCharModel,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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

func GenServiceOptionStreamPageQuery(srvName, methodName, dataModelName, dataModelsPkgPath, facadesPkgPath, gormImplORMPkgPath, servicesPkgCodePath string) error {
	alreadyExistSrvStruct, alreadyExistMethod, codeFilePath, err := checkExistSrvStructOrMethodAndGetCodeFilePath(srvName, methodName, servicesPkgCodePath)
	if err != nil {
		return err
	}

	if alreadyExistMethod {
		return fmt.Errorf("api method: %s.%s already exist", srvName, methodName)
	}

	var (
		code string
		tempNewCodeFilePath string
		needOverrideCodeFile bool
		codeFp *os.File
	)
	if !alreadyExistSrvStruct {
		code, err = getSrvStructCodeContent(srvName, codeFilePath)
		if err != nil {
			return err
		}

		codeFp, err = os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
		if err != nil {
			return err
		}
	} else {
		originalCodeFp, err := os.OpenFile(codeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
		if err != nil {
			return err
		}

		codeBytes, err := ioutil.ReadAll(originalCodeFp)
		if err != nil {
			return err
		}

		code = string(codeBytes) + "\n"

		tempNewCodeFilePath = codeFilePath + ".new"
		codeFp, err = os.OpenFile(tempNewCodeFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
		if err != nil {
			return err
		}

		needOverrideCodeFile = true
	}

	methodContentCode, err := getOptionStreamPageQueryMethodCodeContent(srvName, methodName, dataModelName)
	if err != nil {
		return err
	}

	code += methodContentCode

	if strings.Contains(code, "import (") {
		if !strings.Contains(code, "github.com/995933447/optionstream") {
			code = strings.Replace(code, "import (", `import (
	"github.com/995933447/optionstream"`, 1)
		}

		if !strings.Contains(code, dataModelsPkgPath) {
			code = strings.Replace(code, "import (", `import (
	"` + dataModelsPkgPath + `"`, 1)
		}

		if !strings.Contains(code, facadesPkgPath) {
			code = strings.Replace(code, "import (", `import (
	"` + facadesPkgPath + `"`, 1)
		}

		if !strings.Contains(code, gormImplORMPkgPath) {
			code = strings.Replace(code, "import (", `import (
	"` + gormImplORMPkgPath + `"`, 1)
		}

		if !strings.Contains(code, "context\r\n") {
			code = strings.Replace(code, "import (", `import (
	"context"`, 1)
		}
	}

	outputApiTemplateBytes := len(code)
	var writtenOutputApiTemplateBytes int
	for {
		n, err := codeFp.Write([]byte(code))
		if err != nil {
			return err
		}
		writtenOutputApiTemplateBytes += n
		if writtenOutputApiTemplateBytes >= outputApiTemplateBytes {
			break
		}
		code = code[writtenOutputApiTemplateBytes:]
	}

	if needOverrideCodeFile {
		backupCodeFilePath := codeFilePath + ".backup"
		err = os.Rename(codeFilePath, backupCodeFilePath)
		if err != nil {
			return err
		}
		err = os.Rename(tempNewCodeFilePath, codeFilePath)
		if err != nil {
			return err
		}
		err = os.Remove(backupCodeFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkExistSrvStructOrMethodAndGetCodeFilePath(srvName, methodName, servicesPkgCodePath string) (existSrvStruct, existMethod bool, codeFilePath string, err error) {
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

				if funcDecl, ok := decl.(*ast.FuncDecl); ok &&
					funcDecl.Name.String() == methodName &&
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
					if starExprRecvTypeIdent.String() == methodName {
						existMethod = true
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