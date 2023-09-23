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

package di

import (
	"fmt"
	"strings"
)

const (
	ErrAccessingContainer = "error while accessing DI container"
	ErrAccessingValue     = "error while accessing Value"
	ErrConvertingPointer  = "error while converting pointer"
	ErrNilPointer         = "received nil pointer"
	ErrAssertingType      = "error while asserting type"
)

func ErrGetItemWithKey(key, containerName string, more ...string) error {
	return formatErr(
		ErrAccessingContainer,
		fmt.Sprintf("cannot get item with key %q in container %q", key, containerName),
		formatErr(more...).Error())
}

func formatErr(more ...string) error {
	if len(more) == 0 {
		return fmt.Errorf("")
	}
	return fmt.Errorf(strings.Join(more, ": "))
}
