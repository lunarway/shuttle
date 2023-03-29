package parser

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/kjuulh/shuttletask/pkg/discover"
)

type Function struct {
	Name   string
	Input  []Arg
	Output Output
}

type Arg struct {
	Name string
}

type Output struct {
	Error bool
}

func GenerateAst(ctx context.Context, shuttlelocaldir string, shuttletask *discover.ShuttleTaskDiscovered) ([]*Function, error) {
	funcs := make([]*Function, 0)

	for _, taskfile := range shuttletask.Files {
		tknSet := token.NewFileSet()
		astfile, err := parser.ParseFile(tknSet, path.Join(shuttlelocaldir, "tmp", taskfile), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if ast.FileExports(astfile) {
			decls := astfile.Decls
			for _, decl := range decls {
				funcdecl, ok := decl.(*ast.FuncDecl)
				if ok {
					f := Function{}
					f.Name = funcdecl.Name.Name
					param := funcdecl.Type
					paramList := param.Params.List
					for _, param := range paramList {
						for _, name := range param.Names {
							if name != nil && !strings.Contains(fmt.Sprintf("%s", param.Type), "Context") {
								f.Input = append(f.Input, Arg{
									Name: name.Name,
								})
							}
						}
					}
					outputParam := param.Results
					if outputParam != nil {
						if len(outputParam.List) > 1 {
							return nil, errors.New("only error is supported as an output param")
						}
						if len(outputParam.List) == 0 {
							return nil, errors.New("output params are required, only error is supported")
						}

						for _, param := range outputParam.List {
							if fmt.Sprintf("%s", param.Type) != "error" {
								return nil, errors.New("output was not error")
							}
						}

						f.Output = Output{Error: true}
					}

					funcs = append(funcs, &f)
				}
			}
		}
	}

	return funcs, nil
}
