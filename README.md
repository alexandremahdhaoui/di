# di

Simple dependency injection types for Go.


# How to proceed from now?
- Once we have the list of all ValueFunc/Container Declarations, we can make a graph of their usage.
- What info are we interested with?
  - Is a container built?
  - Is a Value used then mutated?
  - Is a Value used before being Set/Defined?
  - Is Ptr() called on a Value after its container was built?
  - Is Set() called on a Value after its container was built?
  - Is Value() called on a Value before its container was built?

# TODO


di-check:
- Use the `fset.Position(mthd.Pos()).Line` to give users the exact position where a Value is used/called...
  - This will be a useful tool when debugging the dependency graph or to check where a value is used etc..

Immutability:
- Add a "Dirty" field to the di.Value[T]. This information will be used to mark if a Value can possibly be mutated.
  - E.g.: if container is in immature state (not built), the returned values will be marked as Dirty.
  - Therefore, users can create checks on the Value.
- Also, when a Value is taken from a Container, the Container or func returning the Value can mark the Value as "immutable"
and prevent users from calling `Value.Ptr()`.
- Deep copy with DI:
  - Create a field called "Value.Immutable()" returning a bool to let the user know the value is completely immutable.
  - Add options to ValueFunc to enforce immutability when calling the ValueFunc after container build.
    - This means the value must return Dirty equals to false & the underlying type T must implement DeepCopy.

Other todos:
- The only way to ensure immutability is to enforce users to use structs that implements DeepCopy.
- Figure out how to Build a container. This is especially important when consumers of a Value are unaware which
  containers are declared in a Pkg or when producers uses private containers that needs to be built before 
  - A bad pattern would be creating a pointer to the container from the container itself to build the container outside
    the repo.
  - Another way would be to pass the state of the container in the di.Value[T]. But this is also a dirty solution.
    Because Value[T] shouldn't be aware of any constructions related to Container.
  - Finally, Values could implement a IsDirty function which returns a boolean indicator.

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

## Types

| Type      | Description                                                                                                                                                                                                                                                                              |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Container | A container is a data structure holding Values of T.<br/> Container can hold any kind of data. To eventually achieve immutability, the Container implements a Build() function that should be called before accessing data and after setting data.<br/> Container is not a generic type. |
| Value[T]  | Value of T is a generic type used to convey data from and outside a Container                                                                                                                                                                                                            |

## Best effort immutability

Because our system is unable to assert immutability, we make use of `di-check` a binary traversing your code and
constructing a graph of dependencies between the different part of your code that consume and produce Values[T].

## Best practices with DI

In a package, have entry/public interface functions calling ValueFuncs & then internal functions and structs should 
explicitly use the actual T value of Value[T].

Ensuring that DI values are not called everywhere should ease testing & also reduce a bit the complexity. We're using DI
to reduce value propagation and improve cohesion/reduce coupling: so we shouldn't add blow up complexity and add a layer
of hard coupling to DI ValueFunc calls.

## Idea: Dependency injection for distributed systems

To improve cohesion between multiple microservices, we could use dependency injection patterns.

The idea could be to define in your code a `Container` which doesn't refer to an in-memory struct but could be a bridge 
to a "Container service"? or to a shared datastructure such as a k8s ConfigMap.

Another idea, which could be a bit recursive, would be to inject Containers themselves or reference to Containers in
your code, i.e. a user can have an in-memory Container for testing, but when deployed in a kubernetes environment the
user would inject a Distributed Container, which Values would then come from a CM or a container service.

Now this idea brings other issues:
- How do we ensure immutability while we also need this configuration to be changed?
- How do we ensure when your microservice is deployed that the values referenced in the code are available in the 
  Container?
  - The `di-check` could help in that regard. Maybe we could deploy our services using a general easily pluggable 
    operator that checks the microservice dependency graph and check the CM complies to the microservice's specification
