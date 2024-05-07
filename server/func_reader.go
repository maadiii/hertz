package server

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func funcPathAndName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func funcName(f interface{}) string {
	splitFuncName := strings.Split(funcPathAndName(f), ".")

	return splitFuncName[len(splitFuncName)-1]
}

func funcDescription(f interface{}) string {
	fileName, _ := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).FileLine(0)
	funcName := funcName(f)
	fset := token.NewFileSet()

	// Parse src
	parsedAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)

		return ""
	}

	pkg := &ast.Package{
		Name:  "Any",
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = parsedAst
	importPath, _ := filepath.Abs("/")

	myDoc := doc.New(pkg, importPath, doc.AllDecls)
	for _, theFunc := range myDoc.Funcs {
		if theFunc.Name == funcName {
			return theFunc.Doc
		}
	}

	for _, theType := range myDoc.Types {
		for _, theFunc := range theType.Funcs {
			if theFunc.Name == funcName {
				return theFunc.Doc
			}
		}
	}

	return ""
}
