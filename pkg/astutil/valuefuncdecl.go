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
	"go/ast"
	"sync"
)

type ValuefuncDecl struct {
	Decl Decl

	TypeRef      ObjRef
	ContainerRef ObjRef
}

// findFnCalls
// TODO: Implement me
func findFnCalls(fn *ast.FuncDecl) []ObjRef {

	return []ObjRef{{}, {}}
}
func findContainerRef(fn *ast.FuncDecl) (ObjRef, error) {
	return ObjRef{}, nil
}
func findReturnTypes(fn *ast.FuncDecl) []ObjRef {
	sl := make([]ObjRef, 0)

	for _, res := range fn.Type.Results.List {
		index, ok := res.Type.(*ast.IndexExpr)
		if !ok {
			continue
		}

		sel, ok := index.X.(*ast.SelectorExpr)
		if !ok {
			ident, ok := index.X.(*ast.Ident)
			if !ok {
				continue
			}

			sl = append(sl, ObjRef{
				Ident: Ident(ident.Name),
			})
		}

		if sel.Sel == nil {
			continue
		}

		x, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}

		pkgIdent := Ident(x.Name)
		sl = append(sl, ObjRef{
			Ident:    Ident(sel.Sel.Name),
			PkgIdent: &pkgIdent,
		})
	}

	return sl
}

// TODO: Implement me!
func findValueRef(types []ObjRef) (ObjRef, bool) {
	return types[0], true
}

func hasDIFunc(fnCalls []ObjRef) bool {
	for _, fnCall := range fnCalls {
		if fnCall.Ident == "Get" || fnCall.Ident == "Must" || fnCall.Ident == "MustWithOptions" {
			return true
		}
	}

	return false
}

func hasValueType(returnedTypes []ObjRef) bool {
	for _, returnedType := range returnedTypes {
		if returnedType.Ident == "Value" {
			return true
		}
	}

	return false
}

// TODO: Implement
func ValuefuncDeclFromNode(node ast.Node, meta Meta, diPkgIdent Ident) []ValuefuncDecl {
	var wg sync.WaitGroup

	q := make(chan ValuefuncDecl)
	sl := make([]ValuefuncDecl, 0)
	errs := make([]error, 0)

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
		if !hasValueType(returnedTypes) { // inspected func does not return a Value[T] type
			return true
		}

		fnCalls := findFnCalls(fn)
		if len(fnCalls) == 0 { // inspected func did not contain a call to a di function (Get, Must, MustWithOptions...)
			return true
		}
		if !hasDIFunc(fnCalls) {
			return true
		}

		containerRef, err := findContainerRef(fn)
		if err != nil {
			errs = append(errs, err)

			return true
		}

		valueRef, ok := findValueRef(returnedTypes)
		if !ok {
			return true
		}

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
