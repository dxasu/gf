package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/os/gbuild"
)

func main() {
	g.Dump(gbuild.Info())
	g.Dump(gbuild.Map())
}
