package main

import (
	"encoding/json"
	"fmt"

	"github.com/dxasu/gf/frame/g"

	"github.com/dxasu/gf/container/gmap"
)

func main() {
	m := gmap.New()
	m.Sets(g.MapAnyAny{
		"name":  "john",
		"score": 100,
	})
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
