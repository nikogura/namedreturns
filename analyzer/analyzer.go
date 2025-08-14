package analyzer

import (
	"errors"
	"flag"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const FlagReportErrorInDefer = "report-error-in-defer"

var Analyzer = &analysis.Analyzer{
	Name:     "namedreturns",
	Doc:      "Reports functions that don't use named returns",
	Flags:    flags(),
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func flags() flag.FlagSet {
	fs := flag.FlagSet{}
	fs.Bool(FlagReportErrorInDefer, false, "report named error if it is assigned inside defer")
	return fs
}

func run(pass *analysis.Pass) (interface{}, error) {
	reportErrorInDefer := pass.Analyzer.Flags.Lookup(FlagReportErrorInDefer).Value.String() == "true"
	errorType := types.Universe.Lookup("error").Type()

	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("failed to get inspector")
	}

	// only filter function defintions
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		var funcResults *ast.FieldList
		var funcBody *ast.BlockStmt

		switch n := node.(type) {
		case *ast.FuncLit:
			funcResults = n.Type.Results
			funcBody = n.Body
		case *ast.FuncDecl:
			funcResults = n.Type.Results
			funcBody = n.Body
		default:
			return
		}

		// Function without body, ex: https://github.com/golang/go/blob/master/src/internal/syscall/unix/net.go
		if funcBody == nil {
			return
		}

		// no return values - this is fine, no report needed
		if funcResults == nil {
			return
		}

		resultsList := funcResults.List

		// Collect named return variable names
		var namedReturnNames []string
		for _, p := range resultsList {
			if len(p.Names) == 0 {
				// Report this - the parameter is not named and should be
				pass.Reportf(node.Pos(), "unnamed return with type %q found - named returns are required", types.ExprString(p.Type))
				continue
			}

			// Check each name - underscore is not an acceptable return name
			for _, n := range p.Names {
				if n.Name == "_" {
					// Report this - underscore is not a proper name
					pass.Reportf(node.Pos(), "underscore as a return variable name is unacceptable for type %q", types.ExprString(p.Type))
					continue
				}

				// Check if this is an error return that might be exempted
				if !reportErrorInDefer &&
					types.Identical(pass.TypesInfo.TypeOf(p.Type), errorType) &&
					findDeferWithVariableAssignment(funcBody, pass.TypesInfo, pass.TypesInfo.ObjectOf(n)) {
					// This is fine - error return with defer assignment
					continue
				}

				// Collect named return names for later analysis
				namedReturnNames = append(namedReturnNames, n.Name)
			}
		}

		// If we have named returns, check if they're used in return statements and check for shadowing
		if len(namedReturnNames) > 0 {
			checkNamedReturnUsage(pass, funcBody, namedReturnNames, node.Pos())
			checkNamedReturnShadowing(pass, funcBody, namedReturnNames)
		}
	})

	return nil, nil // nolint:nilnil
}

// checkNamedReturnUsage analyzes the function body to see if named return variables are used in return statements
func checkNamedReturnUsage(pass *analysis.Pass, body *ast.BlockStmt, namedReturnNames []string, funcPos token.Pos) {
	ast.Inspect(body, func(node ast.Node) bool {
		if returnStmt, ok := node.(*ast.ReturnStmt); ok {
			// Check if this is a bare return (no expressions)
			if len(returnStmt.Results) == 0 {
				// Bare return is fine when using named returns
				return true
			}

			// Check if the return statement uses the named return variables
			usedNames := make(map[string]bool)
			for _, result := range returnStmt.Results {
				if ident, ok := result.(*ast.Ident); ok {
					// Check if this identifier is one of our named return variables
					for _, namedReturn := range namedReturnNames {
						if ident.Name == namedReturn {
							usedNames[namedReturn] = true
							break
						}
					}
				}
			}

			// Report on named return variables that are declared but not used in this return statement
			for _, namedReturn := range namedReturnNames {
				if !usedNames[namedReturn] {
					pass.Reportf(funcPos, "named return variable %q is declared but not used in return statement", namedReturn)
				}
			}
		}
		return true
	})
}

// checkNamedReturnShadowing detects when named return variables are shadowed by local variables
func checkNamedReturnShadowing(pass *analysis.Pass, body *ast.BlockStmt, namedReturnNames []string) {
	ast.Inspect(body, func(node ast.Node) bool {
		// Check for variable declarations and assignments that might shadow named returns
		switch n := node.(type) {
		case *ast.AssignStmt:
			// Check for := assignments that might shadow named returns
			if n.Tok == token.DEFINE {
				for _, lhs := range n.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						for _, namedReturn := range namedReturnNames {
							if ident.Name == namedReturn {
								pass.Reportf(ident.Pos(), "named return variable %q is shadowed by local variable declaration", namedReturn)
							}
						}
					}
				}
			}
		case *ast.ValueSpec:
			// Check for var declarations that might shadow named returns
			for _, name := range n.Names {
				for _, namedReturn := range namedReturnNames {
					if name.Name == namedReturn {
						pass.Reportf(name.Pos(), "named return variable %q is shadowed by local variable declaration", namedReturn)
					}
				}
			}
		case *ast.RangeStmt:
			// Check for range loop variables that might shadow named returns
			if ident, ok := n.Key.(*ast.Ident); ok {
				for _, namedReturn := range namedReturnNames {
					if ident.Name == namedReturn {
						pass.Reportf(ident.Pos(), "named return variable %q is shadowed by range loop variable", namedReturn)
					}
				}
			}
			if ident, ok := n.Value.(*ast.Ident); ok {
				for _, namedReturn := range namedReturnNames {
					if ident.Name == namedReturn {
						pass.Reportf(ident.Pos(), "named return variable %q is shadowed by range loop variable", namedReturn)
					}
				}
			}
		case *ast.ForStmt:
			// Check for for loop variables that might shadow named returns
			if forStmt, ok := n.Init.(*ast.AssignStmt); ok && forStmt.Tok == token.DEFINE {
				for _, lhs := range forStmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						for _, namedReturn := range namedReturnNames {
							if ident.Name == namedReturn {
								pass.Reportf(ident.Pos(), "named return variable %q is shadowed by for loop variable", namedReturn)
							}
						}
					}
				}
			}
		}
		return true
	})
}

func findDeferWithVariableAssignment(body *ast.BlockStmt, info *types.Info, variable types.Object) bool {
	found := false

	ast.Inspect(body, func(node ast.Node) bool {
		if found {
			return false // stop inspection
		}

		if d, ok := node.(*ast.DeferStmt); ok {
			if fn, ok2 := d.Call.Fun.(*ast.FuncLit); ok2 {
				if findVariableAssignment(fn.Body, info, variable) {
					found = true
					return false
				}
			}
		}

		return true
	})

	return found
}

func findVariableAssignment(body *ast.BlockStmt, info *types.Info, variable types.Object) bool {
	found := false

	ast.Inspect(body, func(node ast.Node) bool {
		if found {
			return false // stop inspection
		}

		if a, ok := node.(*ast.AssignStmt); ok {
			for _, lh := range a.Lhs {
				if i, ok2 := lh.(*ast.Ident); ok2 {
					if info.ObjectOf(i) == variable {
						found = true
						return false
					}
				}
			}
		}

		return true
	})

	return found
}
