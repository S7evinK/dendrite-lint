package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "dendritelint",
	Doc:      "Checks that defined endpoints are used in the internal api server",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.GenDecl)(nil),  // to match const declarations
		(*ast.CallExpr)(nil), // to match method calls
	}

	// Only check package "inthttp"
	if pass.Pkg.Name() != "inthttp" {
		return nil, nil
	}

	knownPaths := make(map[string]ast.Spec)
	usedPaths := make(map[string]bool)

	insp.Preorder(nodeFilter, func(node ast.Node) {

		switch val := node.(type) {
		case *ast.CallExpr:
			// http.Handle has exactly 2 parameters
			if len(val.Args) != 2 {
				return
			}
			fun, ok := val.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			// Currently, only "Handle" methods are checked
			if fun.Sel.Name != "Handle" {
				return
			}
			// "Handle"s first parameter is the path, so we check if it's an
			// identifier (const) which is used. If so, add it to a map.
			ident, ok := val.Args[0].(*ast.Ident)
			if !ok {
				return
			}
			usedPaths[ident.Name] = true
		case *ast.GenDecl:
			// skip everything which is not a const
			if val.Tok != token.CONST {
				return
			}

			for _, spec := range val.Specs {
				for _, vsn := range spec.(*ast.ValueSpec).Names {
					knownPaths[vsn.Name] = spec
				}
			}
		}
	})

	for known, spec := range knownPaths {
		if _, ok := usedPaths[known]; !ok {
			pass.Reportf(spec.Pos(), "declared '%s' endpoint, but not used in internal api server", known)
		}
	}

	return nil, nil
}
