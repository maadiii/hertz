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

func runtimeFunc(f interface{}) *runtime.Func {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer())
}

func funcDescription(f interface{}) string {
	fn := runtimeFunc(f)
	fileName, _ := fn.FileLine(0)
	splitedFuncName := strings.Split(fn.Name(), ".")
	funcName := splitedFuncName[len(splitedFuncName)-1]
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
