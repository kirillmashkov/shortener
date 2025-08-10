// Анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакет main
package exitmainchecker

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// ErrExitMainCheckAnalyzer для использования в анализаторе
var ErrExitMainCheckAnalyzer = &analysis.Analyzer{
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

		ast.Inspect(file, func(node ast.Node) bool {
			mainDecl, isFuncDecl := node.(*ast.FuncDecl)
			if !isFuncDecl {
				return true
			}

			if mainDecl.Name.Name != "main" {
				return false
			}

			ast.Inspect(mainDecl, func(node ast.Node) bool {
				callExpr, isCallExpr := node.(*ast.CallExpr)
				if !isCallExpr {
					return true
				}

				s, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
				if !isSelectorExpr {
					return true
				}

				if s.Sel.Name == "Exit" {
					ident := s.X.(*ast.Ident)
					if ident.Name == "os" {
						pass.Reportf(s.Pos(), "exit call in main function")
					}
				}

				return false
			})

			return false
		})
	}

	return nil, nil
}
