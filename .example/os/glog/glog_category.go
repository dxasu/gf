package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/os/gfile"
	"github.com/dxasu/gf/os/glog"
)

func main() {
	path := "/tmp/glog-cat"
	glog.SetPath(path)
	glog.Stdout(false).Cat("cat1").Cat("cat2").Println("test")
	list, err := gfile.ScanDir(path, "*", true)
	g.Dump(err)
	g.Dump(list)
}
