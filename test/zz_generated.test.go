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

package test

import (
	"github.com/alexandremahdhaoui/di"
	diAstUtil "github.com/alexandremahdhaoui/di/pkg/astutil"
	"go/ast"
)

var (
	MyContainer  = di.New("MyContainer")
	MyContainer2 = di.New("MyContainer2")
)

var MyContainer3 = di.New("MyContainer3")

func MyValueFunc(options ...di.Option) di.Value[map[*diAstUtil.Meta]ast.Node] {
	return di.MustWithOptions[map[*diAstUtil.Meta]ast.Node](MyContainer, "MyValueFunc", options...)
}

func MyValueFunc2(options ...di.Option) di.Value[map[*diAstUtil.Meta]ast.Node] {
	return di.MustWithOptions[map[*diAstUtil.Meta]ast.Node](MyContainer, "MyValueFunc2", options...)
}
