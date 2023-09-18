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
	"fmt"
	//nolint:depguard
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers" //nolint:depguard
)

const ValueFuncSummary = `Creates a single func to conveniently access a di.Value.
This marker is also used by the di-checker to create the dependency graph.
Fields:
	// Name identifies the func that will be used to access the defined value
	Name string
	// Container specifies the Container's Name that will be used to store the Value.
	Container string
	// Type defines the type T to the Value[T].
	Type string
	// TypeImport defines package import for the specific type.
	TypeImport string`

type ValueFunc struct {
	// Name identifies the func that will be used to access the defined value
	Name string
	// Container specifies the Container's Name that will be used to store the Value.
	Container string
	// Type defines the `type` T to the Value[T].
	Type string
	// TypeImport defines package import for the specific type.
	TypeImport string
}

var ValueFuncMarkerDefinition = markers.Must( //nolint:gochecknoglobals
	markers.MakeDefinition("di:valuefunc", markers.DescribesPackage, ValueFunc{}), //nolint:exhaustruct,exhaustivestruct
)

// +controllertools:marker:generateHelp:category="object"

type ValueFuncGenerator struct{}

func (ValueFuncGenerator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, ValueFuncMarkerDefinition); err != nil {
		return err //nolint:wrapcheck
	}

	into.AddHelp(ValueFuncMarkerDefinition, markers.SimpleHelp("object", ValueFuncSummary))

	return nil
}

func (ValueFuncGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		root.NeedTypesInfo()
		markerSet, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
		}

		for _, markerValue := range markerSet[ValueFuncMarkerDefinition.Name] {
			valueFunc := markerValue.(ValueFunc)
			// create the value func with jen
			fmt.Printf("%#v", valueFunc)
		}
	}
	return nil
}
