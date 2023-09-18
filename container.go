package di

import (
	"fmt"
)

const (
	// Options

	InitializeOption Option = iota
	GetValueOption

	// DefaultContainerName is the identifying key to the default Container
	DefaultContainerName = "default"

	// Container State

	immatureContainerState containerState = iota
	builtContainerState
)

var DefaultContainer = New(DefaultContainerName) //nolint:unused

type (
	Container interface {
		get(key string) (pointer, bool)
		set(key string, ptr pointer) error

		// Build an "immutable" container
		Build()
		Name() string
	}

	container struct {
		built    map[string]pointer
		immature map[string]pointer
		name     string
		state    containerState
	}

	pointer *any

	Option         int
	containerState int
)

func (c *container) get(key string) (pointer, bool) {
	if c.state == immatureContainerState {
		v, ok := c.immature[key]

		return v, ok
	}

	v, ok := c.built[key]

	return v, ok
}

func (c *container) set(key string, ptr pointer) error {
	if c.state == builtContainerState {
		return formatErr(fmt.Sprintf("error while setting %q: container %q is immutable", key, c.name))
	}

	c.immature[key] = ptr

	return nil
}

func (c *container) Build() {
	if c.state == builtContainerState {
		return
	}

	for key, ptr := range c.immature {
		newPtr := new(any)
		*newPtr = *ptr
		c.built[key] = newPtr

		delete(c.immature, key)
	}

	c.state = builtContainerState
	c.immature = nil
}

func (c *container) Name() string {
	return c.name
}

func New(name string) Container { //nolint:ireturn
	if name == "" {
		panic("a name is required to create a new container")
	}

	return &container{
		built:    make(map[string]pointer),
		immature: make(map[string]pointer),
		name:     name,
		state:    immatureContainerState,
	}
}

func InitializeValue[T any](c Container, key string) (Value[T], error) {
	v := NewValue[T](key, nil)
	if err := c.set(key, v.pointer()); err != nil {
		return nil, err
	}

	return v, nil
}

func Get[T any](c Container, key string) (Value[T], error) {
	ptr, ok := c.get(key)
	if !ok {
		return nil, ErrGetItemWithKey(key, c.Name())
	}

	converted, err := ConvertPointer[T](ptr)
	if err != nil {
		return nil, err
	}

	return NewValue[T](key, converted), err
}

func Must[T any](c Container, key string) Value[T] {
	v, err := Get[T](c, key)
	if err != nil {
		panic(err)
	}

	return v
}

func MustWithOptions[T any](c Container, key string, option ...Option) Value[T] { //nolint:varnamelen
	// Get -- implicit
	if len(option) == 0 {
		return Must[T](c, key)
	}

	// Get -- explicit
	if option[0] == GetValueOption {
		return Must[T](c, key)
	}

	// Set -- implicit or explicit resolves to Set
	v := NewValue[T](key, nil)
	if err := Set(c, v); err != nil {
		panic(err)
	}

	return v
}

func Set[T any](c Container, v Value[T]) error {
	err := c.set(v.Key(), v.pointer())
	if err != nil {
		return formatErr(fmt.Sprintf("cannot set Value %q to container %q", c.Name(), v.Key()), err.Error())
	}

	return nil
}
