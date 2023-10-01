package astutil

import (
	"encoding/json"
	"fmt"
)

func debug(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b)) //nolint:forbidigo
}

func debugFmt(v interface{}) {
	fmt.Printf("%#v\n", v)
}
