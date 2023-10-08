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
	"fmt"
	"go/format"
	"io"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"strings"
	"unicode"
)

const header = `
//go:build !ignore_autogenerated

%[2]s

// Code generated by di-gen. DO NOT EDIT.
`

func markerName(prefix, name string) string {
	return fmt.Sprintf("%s:%s", prefix, name)
}

func title(s string) string {
	r := []rune(s)

	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}

func generatedFilename(prefix, name string) string {
	return fmt.Sprintf("zz_generated.%s.%s.go", prefix, name)
}

type generateFileOptions struct {
	buffer                           *bytes.Buffer
	ctx                              *genall.GenerationContext
	filename, headerFile, headerYear string
	root                             *loader.Package
}

func generateFile(o generateFileOptions) error {
	var headerText string

	if o.headerFile != "" {
		headerBytes, err := o.ctx.ReadFile(o.headerFile)
		if err != nil {
			return err
		}

		headerText = string(headerBytes)
	}

	headerText = strings.ReplaceAll(headerText, " YEAR", " "+o.headerYear)
	if headerText == "" {
		panic("at the disco!")
	}

	buffer := new(bytes.Buffer)

	_, err := fmt.Fprintf(buffer, header, o.root.Name, headerText)
	if err != nil {
		return err //nolint:wrapcheck
	}

	buffer.Write(o.buffer.Bytes())

	outBytes := buffer.Bytes()
	if formatted, err := format.Source(outBytes); err != nil {
		o.root.AddError(err)
	} else {
		outBytes = formatted
	}

	outputFile, err := o.ctx.Open(o.root, o.filename)
	if err != nil {
		return err //nolint:wrapcheck
	}

	defer func(outputFile io.WriteCloser) {
		err := outputFile.Close()
		if err != nil {
			o.root.AddError(err)
		}
	}(outputFile)

	n, err := outputFile.Write(outBytes)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if n < len(outBytes) {
		return io.ErrShortWrite
	}

	return nil
}