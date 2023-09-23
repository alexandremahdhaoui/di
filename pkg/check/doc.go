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

// Package check
//
// # I. Go Files & Packages
//
// 1. Get current go.mod & module name
//
// 2. List all files with their path & corresponding package
//
// 3. For each file call AST to get list of files where DI pkg is imported & For each file with di pkg imported:
//
// 3.a. List all Container & Get/Must...-function declaration
//
// 3.b. In each file of the pkg: recursively try to find any usage of "locally" declared Container & Func Ident
//
// 4. List Container & ValueFunc Ident in all files of this go module/repo
//
// 5. asda
//
// # Improvements:
//
//   - Create a build function for any Container that is created.
//     That means any created Container (even private ones) can be built (`Container.Build()`)from outside the package.
//
//   - Maybe: recursively parse all go modules used in the GOPATH to check graph, e.g. to check if there is any
//     circular dependencies, or in any other cases where user just want to check all dependencies.
//
// # II. AST
//
// 1. visit AST to get a list of go file where the DI package called + have the ident used in that file
// 2. List declaration of Container (var/const/:=...)
// 3. List declaration of Get/Must/MustWithOptions/Set...-functions
// 4. List Container Ident usage
// 5. List Get/Must...-functions Ident usage/calls
// 6. Make a list of all usages of di.Values[T] and recursively there usages... (have a Max Depth option)
// 7. Make a graph between declaration and usages
//
// # III. Graph
//
// 1. Create a Graph from the list of all vertices (usages of Ident) & Nodes (declaration)
//
// 2. Check if Graph is a DAG (has any circular dependencies)package check
package check
