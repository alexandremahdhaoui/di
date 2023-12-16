# di

Simple dependency injection types for Go.

## Command line tools

In `cmd/`.

2 commands:
- `di-gen`
  - Generate the functions that are used to Get/Set elements in a specific container.
  - Uses the default container by default.

## Types

| Type      | Description                                                                                                                                                                                                                                                                              |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Container | A container is a data structure holding Values of T.<br/> Container can hold any kind of data. To eventually achieve immutability, the Container implements a Build() function that should be called before accessing data and after setting data.<br/> Container is not a generic type. |
| Value[T]  | Value of T is a generic type used to convey data from and outside a Container                                                                                                                                                                                                            |
