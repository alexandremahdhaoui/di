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

// Package computationalgraph
//
// # I. Introduction
//
// Package computationalgraph is used to traverse an AST and build a graph of the simplified computation of the AST.
//
// # II. Usage
//
// Graphs yielded by this package can then be traversed to check e.g. that a di.Value is not used before being set.
//
// # III. Notes
//
//   - We can recursively traverse an AST to build the graph
//
//   - How do we mark/recognize that a func argument has a reference to a di.Value
//
//   - What if a di.Value is an attribute of a struct?
//     We should be able to mark the struct as a holder of a reference to a di.Value
package computationalgraph
