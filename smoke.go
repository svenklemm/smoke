package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {

	var file string
	if len(os.Args) > 1 {
		file = os.Args[1]
	} else {
		file = "smoke.go"
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			checkCall(x)
		}
		return true
	})

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
	fmt.Printf("Call: %s %t\n", func_name, n.Args)
	//fmt.Printf("Function Call\nName: %T\nArgs: %t\n\n", n.Fun, n.Args)

}
