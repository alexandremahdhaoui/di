package astutil

import (
	"go/ast"
)

func exprToObjRef(expr ast.Expr) (ObjRef, bool) {
	switch v := expr.(type) {
	case *ast.Ident:
		return ObjRef{Ident: Ident(v.Name)}, true
	case *ast.SelectorExpr:
		pkgRef, ok := v.X.(*ast.Ident)
		if !ok {
			return ObjRef{}, false
		}

		pkgRefIdent := Ident(pkgRef.Name)

		return ObjRef{
			Ident:    Ident(v.Sel.Name),
			PkgIdent: &pkgRefIdent,
		}, true
	}

	return ObjRef{}, false
}
