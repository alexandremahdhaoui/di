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

package gen

import (
	"bytes"
	"github.com/dave/jennifer/jen"
	"os"
	"path/filepath"

	//nolint:depguard
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers" //nolint:depguard
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../../hack/boilerplate.go.txt,year=2023

// ValueFunc describes a single func that is used to conveniently access a di.Value. This marker is also used by the
// di-checker to create the dependency graph.
type ValueFunc struct {
	// Name identifies the func that will be used to access the defined value
	Name string
	// Container (optional string) specifies the di.Container's Name that will be used to store the Value.
	//    - Container should always resolve to a di.Container defined in the current pkg.
	//    - In other words, the "consumer" of a di.Value, defines both the di.Value and the di.Container in the same
	//      package where the di.Value is consumed.
	//      It's the job of the "producer" of the injectable value to import the ValueFunc from the getter package.
	//      In use cases where an interface is necessary to decouple "consumer" and the "producer", it is a best
	//      practice to create an "interface package" that defines both di.Value & di.Container, which can be imported
	//      by the "consumers" and the "producers" (!! Concurrent producers should NEVER be allowed: greatly reduce the
	//      side effects)
	Container *string
	// Type defines the `type` T to the Value[T].
	Type string
	// TypeImport defines package import for the specific type.
	TypeImport *string
	// Exported indicates if the ValueFunc should be exported or not.
	// Container is exported by default.
	Exported *bool
}

func (vf *ValueFunc) nameWithExportedCasing() string {
	if vf.isExported() {
		return title(vf.Name)
	}

	return vf.Name
}

func (vf *ValueFunc) isExported() bool {
	if vf.Exported == nil {
		return true
	}

	return *vf.Exported
}

var ValueFuncMarkerDefinition = markers.Must( //nolint:gochecknoglobals
	markers.MakeDefinition(markerName(DIMarkerName, ValueFuncMarkerName), markers.DescribesPackage, ValueFunc{}), //nolint:lll,exhaustruct,exhaustivestruct
)

// +controllertools:marker:generateHelp:category="object"

// ValueFuncGenerator Creates a single func to conveniently access a di.Value. This marker is also used by the
// di-checker to create the dependency graph.
//
// Fields:
//
//   - Name (string) identifies the func that will be used to access the defined value.
//
//   - Container (optional string) specifies the di.Container's Name that will be used to store the Value.
//     Container should always resolve to a di.Container defined in the current pkg.
//     In other words, the "consumer" of a di.Value, defines both the di.Value and the di.Container in the same
//     package where the di.Value is consumed.
//     It's the job of the "producer" of the injectable value to import the ValueFunc from the getter package.
//     In use cases where an interface is necessary to decouple "consumer" and the "producer", it is a best
//     practice to create an "interface package" that defines both di.Value & di.Container, which can be imported
//     by the "consumers" and the "producers" (!! Concurrent producers should NEVER be allowed: greatly reduce the
//     side effects)
//
//   - Type (string) defines the type T to the Value[T].
//
//   - TypeImport (optional string) defines package import for the specific type.
//
//   - Exported indicates if the ValueFunc should be exported or not.
//     The ValueFunc is exported by default.
type ValueFuncGenerator struct{}

func (ValueFuncGenerator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, ValueFuncMarkerDefinition); err != nil {
		return err //nolint:wrapcheck
	}

	into.AddHelp(ValueFuncMarkerDefinition, markers.SimpleHelp("object", ""))

	return nil
}

func (ValueFuncGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		root.NeedTypesInfo()

		markerSet, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
		}

		markerValues := markerSet[ValueFuncMarkerDefinition.Name]
		if len(markerValues) == 0 {
			continue
		}

		// We create one zz_generated.di.container.go per package
		// Thus we also instantiate one jen.File per package.
		f := jen.NewFilePath(root.PkgPath) //nolint:varnamelen

		for _, markerValue := range markerValues {
			valueFunc := markerValue.(ValueFunc) //nolint:forcetypeassert

			var typeStt *jen.Statement
			if valueFunc.TypeImport != nil {
				typeStt = jen.Qual(*valueFunc.TypeImport, title(valueFunc.Type))
			} else {
				typeStt = jen.Id(valueFunc.Type)
			}

			var containerStt *jen.Statement
			if valueFunc.Container != nil {
				// We use a jen.Id instead of a jen.Qual because Container should always resolve to a container defined
				// in the current package
				containerStt = jen.Id(*valueFunc.Container)
			} else {
				containerStt = jen.Qual(diPkgPath, "DefaultContainer")
			}

			// func Name(options ...di.Option) di.Value[typeimport.Type] {
			//  	return di.MustWithOptions[typeimport.Type](ContainerName, "Name", options...)
			// }
			f.Func().Id(title(valueFunc.nameWithExportedCasing())).
				Params(jen.Id("options").Op(" ...").Qual(diPkgPath, "Option")).
				Qual(diPkgPath, "Value").Types(typeStt).
				Block(
					jen.Return().Qual(diPkgPath, "MustWithOptions").
						Types(typeStt).
						Call(
							containerStt,
							jen.Lit(valueFunc.nameWithExportedCasing()),
							jen.Id("options").Op("..."),
						),
				)
		}

		buffer := &bytes.Buffer{}
		if err = f.Render(buffer); err != nil {
			return err //nolint:wrapcheck
		}

		filename := filepath.Join(filepath.Dir(root.GoFiles[0]), generatedFilename(DIMarkerName, ValueFuncMarkerName))

		if err = os.WriteFile(filename, buffer.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}

// //+di:valuefunc:name=valueString,type=string

// //+di:valuefunc:container=ContainerName,name=name,type=type,typeImport=github.com/alexandremahdhaoui/type-import
