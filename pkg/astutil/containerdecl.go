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
	"go/token"
)

func varSpecs(node ast.Node) []*ast.ValueSpec {
	sl := make([]*ast.ValueSpec, 0)

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
				sl = append(sl, valueSpec)
			}
		}

		return true
	})

	return sl
}

func ContainerDeclFromNode(node ast.Node, meta Meta) []Decl {
	sl := make([]Decl, 0)

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		for _, spec := range varSpecs(node) {
			// TODO: check the content of the spec.Value[...] to figure out if it's a ContainerDecl or not
			//  Or: Decl owns a pointer to the Values of the declaration

			sl = append(sl, Decl{
				Meta:  meta,
				Ident: Ident(spec.Names[0].Name),
			})
		}

		return true
	})

	return sl
}
