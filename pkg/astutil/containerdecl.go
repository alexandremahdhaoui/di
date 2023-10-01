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
	"go/token"
	"sync"
)

func varSpecs(node ast.Node) []*ast.ValueSpec {
	var wg sync.WaitGroup

	sl := make([]*ast.ValueSpec, 0)
	q := make(chan *ast.ValueSpec)

	go func() {
		for item := range q {
			sl = append(sl, item)

			wg.Done()
		}
	}()

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		if _, ok := node.(*ast.GenDecl); !ok {
			return true
		}

		genDecl := node.(*ast.GenDecl) //nolint:forcetypeassert
		if genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				wg.Add(1)
				q <- valueSpec
			}
		}

		return true
	})

	wg.Wait()
	close(q)

	return sl
}

func ContainerDeclFromNode(node ast.Node, meta Meta, diPkgIdent Ident) []Decl {
	var wg sync.WaitGroup

	sl := make([]Decl, 0)
	q := make(chan Decl)

	go func() {
		for item := range q {
			sl = append(sl, item)

			wg.Done()
		}
	}()

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		for _, spec := range varSpecs(node) {
			for i, variable := range spec.Values {
				callExpr, ok := variable.(*ast.CallExpr)
				if !ok {
					continue
				}

				switch v := callExpr.Fun.(type) {
				// case where diPkgIdent != "." && we have only one var declaration
				case *ast.SelectorExpr:
					if v.Sel.Name != diutil.NewContainerIdent {
						continue
					}

					x, ok := v.X.(*ast.Ident)
					if !ok {
						continue
					}

					// We check if the pkg ident name is the same than the di pkg ident name that was passed to this func.
					if x.Name != diPkgIdent.String() {
						continue
					}
				// case where diPkgIdent == "."
				case *ast.Ident:
					if v.String() != diutil.NewContainerIdent {
						continue
					}
				default:
					continue
				}

				wg.Add(1)
				q <- Decl{
					Meta:  meta,
					Ident: Ident(spec.Names[i].Name),
				}
			}
		}

		return true
	})

	wg.Wait()
	close(q)

	return dropDuplicate(sl)
}

// dropDuplicates takes a slice of Decl, remove any duplicate Ident & return a new slice
func dropDuplicate(sl []Decl) []Decl {
	cleaned := make([]Decl, 0)

	m := make(map[Ident]Decl)
	for _, decl := range sl {
		m[decl.Ident] = decl
	}

	for _, decl := range m {
		cleaned = append(cleaned, decl)
	}

	return cleaned
}
