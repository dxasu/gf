package main

import (
	"fmt"

	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/util/gvalid"
)

func main() {
	g.I18n().SetLanguage("cn")
	err := gvalid.Check("", "required", nil)
	fmt.Println(err.String())
}
