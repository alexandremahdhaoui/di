package computationalgraph

import (
	"github.com/alexandremahdhaoui/di/pkg/astutil"
)

type Assignment struct{}

type StructAttr struct {
}

type Struct struct {
	Ident astutil.Ident

	Attr []StructAttr
}

type FuncArgs struct{}

type FuncCall struct {
	Meta astutil.Meta

	FuncIdent astutil.Ident
	Args      []FuncArgs
}

type AssociatedFuncCall struct {
	Meta astutil.Meta

	Args                []FuncArgs
	AssociatedFuncIdent astutil.Ident
}

type IfBranch struct {
	Meta astutil.Meta
}

type SwitchBranch struct {
	Meta astutil.Meta
}

type ForLoop struct {
	Meta astutil.Meta
}
