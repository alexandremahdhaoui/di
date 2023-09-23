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
	"fmt"
	"github.com/alexandremahdhaoui/di/pkg/astutil"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

func debug(v interface{}) {
	fmt.Printf("%v\n", v) //nolint:forbidigo
}

func VisitAST() {
	roots, err := loader.LoadRoots("./test") // use "..." to load all packages
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
			cSl := astutil.ContainerDeclFromNode(node, meta)
			vSl := astutil.ValuefuncDeclFromNode(node, meta)

			if len(pkgImports) > 0 {
				debug(pkgImports)
			}

			if len(vSl) > 0 {
				debug(vSl)
			}

			if len(cSl) > 0 {
				debug(cSl)
			}
		}
	}
}
