# di

Simple dependency injection types for Go.

# TODO

- The only way to ensure immutability is to enforce users to use structs that implements DeepCopy.
- Figure out how to Build a container. This is especially important when consumers of a Value are unaware which
  containers are declared in a Pkg or when producers uses private containers that needs to be built before 
  - A bad pattern would be creating a pointer to the container from the container itself to build the container outside
    the repo.
  - Another way would be to pass the state of the container in the di.Value[T]. But this is also a dirty solution.
    Because Value[T] shouldn't be aware of any constructions related to Container.
  - Finaly, Values could implement a IsDirty function which returns a boolean indicator.

## Command line tools

In `cmd/`.

2 commands:
- `di-check`
  - Build a dependency graph & check:
    - No values are called if not Set or Defined upfront
    - No values are being set after its container was built
    - No values are being accessed if its container was not built upfront
- `di-gen`
  - Generate the functions that are used to Query/Set elements in a specific container.
  - Uses the default container by default.
