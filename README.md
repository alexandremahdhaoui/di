# di

Simple dependency injection types for Go.

## Generator

In `cmd/di`.

2 sub-commands:
- `check`
  - Build a dependency graph & check:
    - No values are called if not Set or Defined upfront
    - No values are being set after its container was built
    - No values are being accessed if its container was not built upfront
- `generate`
  - Generate the functions that are used to Query/Set elements in a specific container.
  - Uses the default container by default.
