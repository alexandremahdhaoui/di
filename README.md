# di

Simple dependency injection types for Go.

# TODO

- Migrate pkg/gen -> valuefunc & container Generate() func to use genall internals [genall/output.go#L145](/var/home/alex/go/pkg/mod/sigs.k8s.io/controller-tools@v0.13.0/pkg/genall/output.go)
- Create di-check binary: visit ast to build the dependency graph
  - https://astexplorer.net/ 
  - https://yuroyoro.github.io/goast-viewer/
  - https://www.zupzup.org/go-ast-traversal/index.html

## CLIs

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
