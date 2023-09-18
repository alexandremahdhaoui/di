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
