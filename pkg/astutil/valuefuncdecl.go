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
	"github.com/alexandremahdhaoui/di/pkg/diutil"
	"go/ast"
	"sync"
)

type ValuefuncDecl struct {
	Decl Decl

	TypeRef      ObjRef
	ContainerRef ObjRef
}

func findContainerRef(fn *ast.FuncDecl, diPkgIdent Ident) (ObjRef, bool) {
	for _, v := range fn.Body.List {
		switch stmt := v.(type) {
		case *ast.ReturnStmt:
			for _, res := range stmt.Results {
				switch expr := res.(type) {
				case *ast.CallExpr:
					// We check if the called fun is a ref to a di function that yields a Value (so Get, Must...)
					indexExpr, ok := expr.Fun.(*ast.IndexExpr)
					if !ok {
						continue
					}

					switch expr := indexExpr.X.(type) {
					// Case where DI is normally imported
					case *ast.SelectorExpr:
						objRef, ok := exprToObjRef(expr)
						if !ok {
							continue
						}

						if objRef.PkgIdent == nil {
							// We move on if diPkgIdent is not a dot import
							if diPkgIdent.String() == diutil.DotImportIdent {
								continue
							}
						} else {
							// We move on if the pkg does not refer to DI
							if objRef.PkgIdent.String() != diPkgIdent.String() {
								continue
							}
						}

						fn := objRef.Ident.String()
						// we move on if the function does not retrieve a value.
						if fn != diutil.GetIdent && fn != diutil.MustIdent && fn != diutil.MustWithOptionsIdent {
							continue
						}
					case *ast.Ident:
						if expr.String() != diPkgIdent.String() {
							continue
						}
					default:
						continue
					}

					// ContainerRef is the first arg in the func call
					return exprToObjRef(expr.Args[0])
				default:
					continue
				}
			}
		default:
			continue
		}
	}
	return ObjRef{}, false
}

func findReturnTypes(fn *ast.FuncDecl) []ObjRef {
	sl := make([]ObjRef, 0)

	for _, res := range fn.Type.Results.List {
		index, ok := res.Type.(*ast.IndexExpr)
		if !ok {
			continue
		}

		objRef, ok := exprToObjRef(index.X)
		if !ok {
			continue
		}

		sl = append(sl, objRef)
	}

	return sl
}

func findValueRef(refs []ObjRef, diPkgIdent Ident) (ObjRef, bool) {
	for _, ref := range refs {
		if ref.Ident != diutil.ValueIdent {
			return ObjRef{}, false
		}

		if ref.PkgIdent == nil {
			if diPkgIdent.String() == diutil.DotImportIdent {
				return ref, true
			}
		}

		if ref.PkgIdent.String() == diPkgIdent.String() {
			return ref, true
		}
	}

	return ObjRef{}, false
}

func hasValueType(refs []ObjRef, diPkgIdent Ident) bool {
	for _, ref := range refs {
		if ref.Ident == diutil.ValueIdent {
			// Pkg reference ident is nil, so we'll check if DI is dot import
			if ref.PkgIdent == nil {
				if diPkgIdent.String() == diutil.DotImportIdent {
					return true
				}
			}
			// We check if the pkg ref ident matches
			if ref.PkgIdent.String() == diPkgIdent.String() {
				return true
			}
		}
	}

	return false
}

func ValuefuncDeclFromNode(node ast.Node, meta Meta, diPkgIdent Ident) []ValuefuncDecl {
	var wg sync.WaitGroup

	q := make(chan ValuefuncDecl)
	sl := make([]ValuefuncDecl, 0)

	go func() {
		for item := range q {
			sl = append(sl, item)

			wg.Done()
		}
	}()

	ast.Inspect(node, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncDecl); !ok {
			return true
		}

		fn := node.(*ast.FuncDecl) //nolint:forcetypeassert

		// return early if func doesn't return any types
		// returnedTypes must have a Value[T]
		returnedTypes := findReturnTypes(fn)
		if len(returnedTypes) == 0 {
			return true
		}
		if !hasValueType(returnedTypes, diPkgIdent) { // inspected func does not return a Value[T] type
			return true
		}

		containerRef, ok := findContainerRef(fn, diPkgIdent)
		if !ok {
			return true
		}

		valueRef, ok := findValueRef(returnedTypes, diPkgIdent)
		if !ok {
			return true
		}

		wg.Add(1)
		q <- ValuefuncDecl{
			Decl: Decl{
				Meta:  meta,
				Ident: Ident(fn.Name.Name),
			},
			TypeRef:      valueRef,
			ContainerRef: containerRef,
		}

		return true
	})

	wg.Wait()
	close(q)

	return sl
}
