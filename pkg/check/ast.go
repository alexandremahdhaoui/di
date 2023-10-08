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

package check

import (
	"encoding/json"
	"fmt"
	"github.com/alexandremahdhaoui/di/pkg/astutil"
	"github.com/alexandremahdhaoui/di/pkg/diutil"
	"github.com/alexandremahdhaoui/graph"
	"path/filepath"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

func debug(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b)) //nolint:forbidigo
}

func DiPkgIdent(pkgImports []astutil.PkgImport) (astutil.Ident, bool) {
	for _, pkgImport := range pkgImports {
		if pkgImport.Pkg() == diutil.PkgPath {
			return pkgImport.Ident(), true
		}
	}

	return "", false
}

func VisitAST(path string) {
	localPackages := LocalPackages("./...")
	debug(localPackages)

	roots, err := loader.LoadRoots(path) // use "..." to load all packages
	if err != nil {
		panic(err)
	}

	for _, root := range roots {
		root.NeedTypesInfo()
		root.NeedSyntax()

		for i, file := range root.GoFiles {
			node := root.Syntax[i]
			meta := astutil.Meta{
				Pkg:      root.Package.Name,
				Filepath: file,
				Module:   root.PkgPath,
			}

			pkgImports := astutil.PkgImportFromNode(node)

			diPkgIdent, ok := DiPkgIdent(pkgImports)
			if !ok {
				continue
			}

			cSl := astutil.ContainerDeclFromNode(node, meta, diPkgIdent)
			vSl := astutil.ValuefuncDeclFromNode(node, meta, diPkgIdent)

			if len(vSl) > 0 {
				debug(vSl)
			}

			if len(cSl) > 0 {
				debug(cSl)
			}
		}
	}
}

// LocalPackages returns a slice of available pkg definitions from the given []*loader.Package
//
// Usage:
//
//   - These LocalPackages are used when recursively traversing the tree of the given function (e.g. a main func) to
//     know if a package is locally defined, and therefore subject to be traversed.
//
//   - If we encounter a function call with a local package ref, we will traverse that function and recursively perform
//     the same operations on its sequence of statement.
//     This operation will build a func relationship graph.
func LocalPackages(path string) []astutil.Meta {
	sl := make([]astutil.Meta, 0)

	roots, err := loader.LoadRoots(path) // use "..." to load all packages
	if err != nil {
		panic(err)
	}

	for _, root := range roots {
		sl = append(sl, astutil.Meta{
			Pkg:      root.Package.Name,
			Filepath: filepath.Dir(root.GoFiles[0]),
			Module:   root.PkgPath,
		})
	}

	return sl
}

// FuncRelationshipGraph traverses a given function and finds statement and expressions that refers to other function
// calls in order to produce a slice of nodes and vertices.
//
// Each function is a node & each function call in a func is one of its vertices.
// Vertices are ordered, by order of apparition in the func's body.
//
// Usage:
//
//	We will traverse that tree/graph to find if a di.Value was used before being Set, etc...
func FuncRelationshipGraph() *graph.Node {

}

// DIRelationshipGraph traverses a given function to find any reference to a DI ValueFunc, di.Container or to a DI
// Value.
// We will also track Ident which have a di ValueFunc, di.Container or di.Value
func DIRelationshipGraph() {}
