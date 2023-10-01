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
	"sync"
)

type PkgImport struct {
	name *string
	path string
}

// Pkg returns the package path, e.g.: "github.com/alexandremahdhaoui/di"
func (d *PkgImport) Pkg() string {
	return strings.ReplaceAll(d.path, "\"", "")
}

// Ident returns the user specified ident for the package or the package's name if no user specified ident was found
func (d *PkgImport) Ident() Ident {
	if d.name != nil {
		return Ident(*d.name)
	}

	sl := strings.Split(d.Pkg(), "/")

	return Ident(sl[len(sl)-1])
}

func NewPkgImport(spec *ast.ImportSpec) PkgImport {
	var importName *string

	if spec.Name != nil {
		name := spec.Name.Name
		importName = &name
	}

	if spec.Path == nil {
		return PkgImport{}
	}

	return PkgImport{
		name: importName,
		path: spec.Path.Value,
	}
}

func importSpecs(node ast.Node) []*ast.ImportSpec {
	var wg sync.WaitGroup

	q := make(chan *ast.ImportSpec)
	sl := make([]*ast.ImportSpec, 0)

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

		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		if genDecl.Tok != token.IMPORT {
			return true
		}

		for _, spec := range genDecl.Specs {
			if importSpec, ok := spec.(*ast.ImportSpec); ok {
				wg.Add(1)
				q <- importSpec
			}
		}

		return true
	})

	wg.Wait()
	close(q)

	return sl
}

func PkgImportFromNode(node ast.Node) []PkgImport {
	var wg sync.WaitGroup

	q := make(chan PkgImport)
	sl := make([]PkgImport, 0)

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

		for _, spec := range importSpecs(node) {
			wg.Add(1)
			q <- NewPkgImport(spec)
		}

		return true
	})

	wg.Wait()
	close(q)

	return sl
}
