package main

import (
	"fmt"

	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/os/gres"
	_ "github.com/dxasu/gf/os/gres/testdata"
)

func main() {
	gres.Dump()

	v := g.View()
	s, err := v.Parse("index.html")
	fmt.Println(err)
	fmt.Println(s)
}
