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
	"unsafe"
)

type (
	Value[T any] interface {
		pointer() pointer
		Key() string

		MustPtr() *T
		Ptr() (*T, error)

		MustSet(item T)
		Set(item T) error

		MustValue() T
		Value() (T, error)
	}

	value[T any] struct {
		key string
		ptr pointer
	}
)

func (v *value[T]) pointer() pointer {
	return v.ptr
}

func (v *value[T]) Key() string {
	return v.key
}

func (v *value[T]) MustPtr() *T {
	ptr, err := v.Ptr()
	if err != nil {
		panic(err)
	}

	return ptr
}

func (v *value[T]) Ptr() (*T, error) {
	return ConvertPointer[T](v.ptr)
}

func (v *value[T]) MustSet(item T) {
	if err := v.Set(item); err != nil {
		panic(err)
	}
}

func (v *value[T]) Set(item T) error {
	*v.ptr = any(item)

	return nil
}

func (v *value[T]) MustValue() T { //nolint:ireturn
	val, err := v.Value()
	if err != nil {
		panic(err)
	}

	return val
}

func (v *value[T]) Value() (T, error) { //nolint:ireturn
	if v.ptr == nil {
		return *new(T), formatErr("cannot return value", ErrNilPointer)
	}

	ptr, err := v.Ptr()
	if err != nil {
		return *new(T), err
	}

	return *ptr, nil

}

func ConvertPointer[T any](ptr pointer) (*T, error) {
	if ptr == nil {
		return nil, formatErr(ErrConvertingPointer, ErrNilPointer)
	}

	// deref := *ptr
	// // try to type assert the concrete data pointer by ptr
	// if _, ok := deref.(T); !ok {
	// 	return nil, formatErr(ErrConvertingPointer, ErrAssertingType)
	// }

	return (*T)(unsafe.Pointer(ptr)), nil
}

func NewValue[T any](key string, pointer *T) Value[T] {
	if key == "" {
		panic("a key is required to create a new Value")
	}

	if pointer == nil {
		return &value[T]{
			key: key,
			ptr: (*any)(unsafe.Pointer(new(T))),
		}
	}

	return &value[T]{
		key: key,
		ptr: (*any)(unsafe.Pointer(pointer)),
	}
}
