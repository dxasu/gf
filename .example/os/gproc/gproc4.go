package main

import (
	"os"
	"time"

	"github.com/dxasu/gf/os/genv"
	"github.com/dxasu/gf/os/glog"
	"github.com/dxasu/gf/os/gproc"
)

// 查看父子进程的环境变量
func main() {
	time.Sleep(5 * time.Second)
	glog.Printf("%d: %v", gproc.Pid(), genv.All())
	p := gproc.NewProcess(os.Args[0], os.Args, os.Environ())
	p.Start()
}
