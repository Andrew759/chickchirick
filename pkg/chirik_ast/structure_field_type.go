package chirik_ast

import (
	"fmt"
	"go/ast"
	"strings"
)

type Type struct {
	ast ast.Expr
}

func (t *Type) String() string {
	switch t.ast.(type) {
	case *ast.Ident:
		ident := t.ast.(*ast.Ident)

		return ident.Name
	case *ast.SelectorExpr:
		selectorExpr := t.ast.(*ast.SelectorExpr)

		return fmt.Sprintf("%s.%s", selectorExpr.X, selectorExpr.Sel)
	case *ast.StarExpr:
		starExpr := t.ast.(*ast.StarExpr)

		if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
			return fmt.Sprintf("*%s.%s", selectorExpr.X, selectorExpr.Sel)
		} else {
			astExpr := starExpr.X.(*ast.Ident)
			return fmt.Sprintf("*%s", astExpr)
		}
	default:
		panic("unexpected struct field type")
	}
}

func (t *Type) Value() string {
	return strings.TrimPrefix(t.String(), "*")
}

func (t *Type) IsNullable() bool {
	return strings.HasPrefix(t.String(), "*")
}
