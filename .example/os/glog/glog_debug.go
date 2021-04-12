package main

import (
	"time"

	"github.com/dxasu/gf/os/glog"
	"github.com/dxasu/gf/os/gtime"
	"github.com/dxasu/gf/os/gtimer"
)

func main() {
	gtimer.SetTimeout(3*time.Second, func() {
		glog.SetDebug(false)
	})
	for {
		glog.Debug(gtime.Datetime())
		time.Sleep(time.Second)
	}
}
