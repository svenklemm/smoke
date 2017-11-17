package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
	"strings"
)

var targets = []string{"fmt.Printf"}

func main() {

	// dir := "pkg/snippets"
	//if len(os.Args) > 1 {
	//	dir = os.Args[1]
	//}

	fset := token.NewFileSet() // positions are relative to fset
	/*
		f, err := parser.ParseDir(fset, dir, nil, 0)
		if err != nil {
			panic(err)
		}

		for pkg, tree := range f {
			fmt.Printf("\nPACKAGE: %s\n\n", pkg)
			ast.Inspect(tree, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					checkCall(x)
				}
				return true
			})
		}
	*/

	f, err := parser.ParseFile(fset, "pkg/snippets/snippets.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}
	files := []*ast.File{f}

	// Create the type-checker's package.
	pkg := types.NewPackage("snippets", "")

	// Type-check the package, load dependencies.
	// Create and build the SSA program.
	hello, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()}, fset, pkg, files, ssa.SanityCheckFunctions)
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Print out the package.
	hello.WriteTo(os.Stdout)

	// Print out the package-level functions.
	hello.Func("foobar").WriteTo(os.Stdout)

}

func extractName(n *ast.SelectorExpr) []string {
	var fun []string
	switch x := n.X.(type) {
	case *ast.Ident:
		fun = append(fun, x.String())
	case *ast.SelectorExpr:
		fun = extractName(x)
	default:
		fmt.Printf("UNKNOWN NODE: %T", x)
	}
	fun = append(fun, n.Sel.String())
	return fun
}

func getFunctionName(n ast.Expr) string {
	var fun []string
	switch name := n.(type) {
	case *ast.Ident:
		fun = append(fun, name.String())
	case *ast.SelectorExpr:
		fun = extractName(name)
	default:
		fmt.Printf("UNKNOWN NODE: %T", name)
	}
	return strings.Join(fun, ".")
}

func checkCall(n *ast.CallExpr) {
	func_name := getFunctionName(n.Fun)

	for _, f := range targets {
		if f == func_name {
			fmt.Printf("Call: %s %t\n", func_name, n.Args)
		} else {
			fmt.Printf("Call: %s %t\n", func_name, n.Args)
		}
	}
}
