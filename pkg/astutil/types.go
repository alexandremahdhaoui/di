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
	"path/filepath"
)

type (
	Ident string

	Meta struct {
		// Pkg is the package name where the object is defined
		Pkg string
		// Filepath where the object is defined
		Filepath string
		// Module in which the object is defined
		Module string
		// Pos where the token is defined
		Pos token.Pos
	}

	ObjRef struct {
		// Ident of the object we are referencing
		Ident Ident
		// PkgIdent (optional) is the identity referencing a PkgImport
		PkgIdent *Ident
	}

	// Decl example: var Container = di.New("container")
	Decl struct {
		Meta  Meta
		Ident Ident
	}

	Usage struct {
		Meta   Meta
		ObjRef ObjRef
	}
)

func (m *Meta) PkgPath() string {
	return filepath.Dir(m.Filepath)
}

func (i Ident) Exported() bool {
	return ast.IsExported(string(i))
}

func (i Ident) String() string {
	return string(i)
}
