// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package findcall defines an Analyzer that serves as a trivial
// example and test of the Analysis API. It reports a diagnostic for
// every call to a function or method of the name specified by its
// -name flag. It also exports a fact for each declaration that
// matches the name, plus a package-level fact if the package contained
// one or more such declarations.
package valerror

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const Doc = `find calls to a particular function

The findcall analysis reports calls to functions or methods
of a particular name.`

var Analyzer = &analysis.Analyzer{
	Name:             "error_linter",
	Doc:              Doc,
	URL:              "https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/findcall",
	Run:              run,
	RunDespiteErrors: true,
	FactTypes:        []analysis.Fact{},
}

var name string // -name flag

func init() {
	Analyzer.Flags.StringVar(&name, "name", name, "name of the function to find")
}
func run(pass *analysis.Pass) (interface{}, error) {
	name := "|" + strings.ReplaceAll(name, ",", "|") + "|"
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				if _, ok := pass.TypesInfo.Defs[decl.Name].(*types.Func); ok {
					ast.Inspect(decl, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							var id *ast.Ident
							funcName := ""
							switch fun := call.Fun.(type) {
							case *ast.Ident:
								id = fun
								funcName = fun.Name
							case *ast.SelectorExpr:
								funcName = fun.Sel.Name
								id = fun.Sel
								if fun.X != nil {
									if tmp, ok := fun.X.(*ast.Ident); ok {
										funcName = tmp.Name + "." + funcName
									}
								}
							}

							if id != nil && !pass.TypesInfo.Types[id].IsType() && strings.Contains(name, "|"+funcName+"|") {
								pass.Report(analysis.Diagnostic{
									Pos:     call.Lparen,
									Message: fmt.Sprintf("call of %s(...)", funcName),
								})
							}
						}
						return true
					})

				}
			}
		}
	}

	return nil, nil
}
