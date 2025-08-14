// Анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакет main
package exitmainchecker

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer для использования в анализаторе
var Analyzer = &analysis.Analyzer{
	Name: "exitmaincheck",
	Doc:  "check call exit from main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if fullPath := pass.Fset.Position(file.Pos()).String(); strings.Contains(fullPath, "go-build") {
			continue
		}
		if pass.Pkg.Name() != "main" {
			continue
		}

		var inMainFunc bool
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				inMainFunc = n.Name.Name == "main"
				return true
			case *ast.CallExpr:
				if inMainFunc {
					s, isSelectorExpr := n.Fun.(*ast.SelectorExpr)
					if isSelectorExpr && s.Sel.Name == "Exit" {
						ident, isIdent := s.X.(*ast.Ident)
						if isIdent && ident.Name == "os" {
							pass.Reportf(s.Pos(), "exit call in main function")
						}
					}
				}
				return true
			}
			return true
		})
	}
	return nil, nil
}
