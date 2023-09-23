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
	"strings"
)

type PkgImport struct {
	name *string
	path string
}

func (d *PkgImport) Pkg() string {
	return strings.ReplaceAll(d.path, "\"", "")
}

func (d *PkgImport) Ident() Ident {
	if d.name != nil {
		return Ident(*d.name)
	}

	sl := strings.Split(d.Pkg(), "/")

	return Ident(sl[len(sl)-1])
}

func NewPkgImport(spec *ast.ImportSpec) *PkgImport {
	var importName *string

	if spec.Name != nil {
		name := spec.Name.Name
		importName = &name
	}

	if spec.Path == nil {
		return nil
	}

	return &PkgImport{
		name: importName,
		path: spec.Path.Value,
	}
}

func importSpecs(node ast.Node) []*ast.ImportSpec {
	sl := make([]*ast.ImportSpec, 0)

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		if _, ok := node.(*ast.GenDecl); !ok {
			return true
		}

		genDecl := node.(*ast.GenDecl) //nolint:forcetypeassert
		if genDecl.Tok != token.IMPORT {
			return true
		}

		for _, spec := range genDecl.Specs {
			if importSpec, ok := spec.(*ast.ImportSpec); ok {
				sl = append(sl, importSpec) //nolint:forcetypeassert
			}

		}

		return true
	})

	return sl
}

func PkgImportFromNode(node ast.Node) []PkgImport {
	sl := make([]PkgImport, 0)

	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		for _, spec := range importSpecs(node) {
			var name *string

			if spec.Name != nil {
				s := spec.Name.Name
				name = &s
			}

			sl = append(sl, PkgImport{
				name: name,
				path: spec.Path.Value,
			})
		}

		return true
	})

	return sl
}
