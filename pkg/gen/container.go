package gen

import (
	"fmt"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

const ContainerSummary = `di:container marker conveniently generates a new Container`

// Container (`di:container` marker) conveniently generates a new Container
type Container struct {
	// Name identifies the container to be created
	Name string
	// IsPrivate indicates if the Container should be exported or not.
	// Container is public by default, i.e. IsPrivate is false by default.
	IsPrivate bool
}

var ContainerMarkerDefinition = markers.Must( //nolint:gochecknoglobals
	markers.MakeDefinition("di:container", markers.DescribesPackage, Container{}), //nolint:exhaustruct,exhaustivestruct
)

// +controllertools:marker:generateHelp:category="object"

type ContainerGenerator struct{}

func (ContainerGenerator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, ContainerMarkerDefinition); err != nil {
		return err //nolint:wrapcheck
	}

	into.AddHelp(ContainerMarkerDefinition, markers.SimpleHelp("object", ContainerSummary))

	return nil
}

func (ContainerGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		root.NeedTypesInfo()
		markerSet, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
		}

		for _, markerValue := range markerSet[ContainerMarkerDefinition.Name] {
			container := markerValue.(Container)
			// do stuff with the marker, e.g. creates a var declaration to create the new container
			fmt.Printf("%#v", container)
		}
	}

	return nil
}
