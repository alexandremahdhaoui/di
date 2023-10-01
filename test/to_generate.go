package test

//+di:container:name=testContainer0,exported=true

//+di:container:name=testContainer1,exported=false

//+di:valuefunc:name=testValueFunc0,type=string

//+di:valuefunc:container=TestContainer0,name=testValueFunc1,type=di.Value[map[*astutil.Meta]ast.Node],typeImport=github.com/alexandremahdhaoui/di/pkg/astutil
