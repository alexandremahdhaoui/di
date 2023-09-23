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

//nolint:depguard
import (
	"bytes"
	"github.com/dave/jennifer/jen"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../../hack/boilerplate.go.txt,year=2023

type Container struct {
	// Name identifies the container to be created
	Name string
	// Exported indicates if the Container should be exported or not.
	// The Container is not exported by default.
	Exported *bool
}

func (c *Container) nameWithExportedCasing() string {
	if c.isExported() {
		return title(c.Name)
	}

	return c.Name
}

func (c *Container) isExported() bool {
	if c.Exported == nil {
		return false
	}

	return *c.Exported
}

var ContainerMarkerDefinition = markers.Must( //nolint:gochecknoglobals
	markers.MakeDefinition(markerName(DIMarkerName, ContainerMarkerName),
		markers.DescribesPackage, Container{}), //exhaustruct,exhaustivestruct
)

// +controllertools:marker:generateHelp:category="object"

// ContainerGenerator Conveniently generates a new Container
//
// Fields:
//
//   - Name (string) identifies the container to be created.
//
//   - Exported (optional bool) indicates if the Container should be exported or not.
//     The Container is not exported by default.
type ContainerGenerator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
}

func (ContainerGenerator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, ContainerMarkerDefinition); err != nil {
		return err //nolint:wrapcheck
	}

	into.AddHelp(ContainerMarkerDefinition, markers.SimpleHelp("object", ""))

	return nil
}

func (g ContainerGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		root.NeedTypesInfo()

		markerSet, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
		}

		markerValues := markerSet[ContainerMarkerDefinition.Name]
		if len(markerValues) == 0 {
			continue
		}

		// We create one zz_generated.di.container.go per package
		// Thus we also instantiate one jen.File per package.
		f := jen.NewFilePath(root.PkgPath) //nolint:varnamelen
		varDefinitions := make([]jen.Code, 0)

		for _, markerValue := range markerValues {
			container := markerValue.(Container) //nolint:forcetypeassert

			varDefinitions = append(varDefinitions, jen.
				Id(container.nameWithExportedCasing()).
				Op("=").
				Qual(diPkgPath, "New").
				Call(jen.Lit(container.nameWithExportedCasing())))
		}

		f.Var().Defs(varDefinitions...)

		buffer := &bytes.Buffer{}
		if err = f.Render(buffer); err != nil {
			root.AddError(err)

			return err
		}

		if err = generateFile(generateFileOptions{
			buffer:     buffer,
			ctx:        ctx,
			filename:   generatedFilename(DIMarkerName, ContainerMarkerName),
			headerFile: g.HeaderFile,
			headerYear: g.Year,
			root:       root,
		}); err != nil {
			root.AddError(err)

			return err
		}
	}

	return nil
}

//+di:container:name=container0,exported=true

//+di:container:name=container1,exported=false
