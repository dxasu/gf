package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/os/glog"
)

func main() {
	err := glog.SetConfigWithMap(g.Map{
		"prefix": "[TEST]",
	})
	if err != nil {
		panic(err)
	}
	glog.Info(1)
}
