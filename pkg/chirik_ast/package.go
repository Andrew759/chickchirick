package chirik_ast

import "go/ast"

type Package struct {
	ast *ast.Ident
}

func (p *Package) Name() string {
	return p.ast.Name
}

func (p *Package) SetName(name string) {
	p.ast.Name = name
}
