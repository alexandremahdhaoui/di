/*
Copyright 2023 Alexandre Mahdhaoui.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package astutil

import (
	"fmt"
	"go/ast"
)

type ValuefuncDecl struct {
	Decl Decl

	TypeRef      ObjRef
	ContainerRef ObjRef
}

func debug(v interface{}) {
	fmt.Printf("%v\n", v) //nolint:forbidigo
}

func ValuefuncDeclFromNode(node ast.Node, meta Meta) []ValuefuncDecl {
	sl := make([]ValuefuncDecl, 0)

	ast.Inspect(node, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncDecl); !ok {
			return true
		}

		fn := node.(*ast.FuncDecl)

		// find a way to recursively parse type of
		//  - fn.Type.Results.List[...]
		//  - fn.Body.List[...] (basically we want to access something like `.List[-1]` to check the return stmt

		debug(fn)

		// With the content of what is being parsed below, we should be able to parse all generated functions
		// This can then be used to construct a graph of dependencies between consumer and producer code
		fmt.Println("---\nfound a func declaration")

		debug("\t---")
		fmt.Printf("\t%#v\n", fn)

		debug("\t---")
		fmt.Printf("\t%#v\n", fn.Type)
		fmt.Printf("\t%#v\n", fn.Type.Results)
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).X)
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel)
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index)
		//fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.SelectorExpr).X)
		//fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.SelectorExpr).Sel)
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.MapType).Key)
		// We can check star references
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.MapType).Key.(*ast.StarExpr))
		fmt.Printf("\t%#v\n", fn.Type.Results.List[0].Type.(*ast.IndexExpr).Index.(*ast.MapType).Value)

		debug("\t---")
		fmt.Printf("\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Fun.(*ast.IndexExpr).X.(*ast.SelectorExpr))
		fmt.Printf("\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Fun.(*ast.IndexExpr).X.(*ast.SelectorExpr).X)
		fmt.Printf("\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Fun.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel)
		fmt.Printf("\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Fun.(*ast.IndexExpr).Index)

		debug("\t\t---")
		fmt.Printf("\t\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Args[0])
		fmt.Printf("\t\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Args[1])
		fmt.Printf("\t\t%#v\n", fn.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CallExpr).Args[2])

		return true
	})

	return sl
}
