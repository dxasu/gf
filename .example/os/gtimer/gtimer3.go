package main

import (
	"time"

	"github.com/dxasu/gf/os/glog"
	"github.com/dxasu/gf/os/gtimer"
)

func main() {
	interval := time.Second
	gtimer.AddSingleton(interval, func() {
		glog.Println("doing")
		time.Sleep(5 * time.Second)
	})

	select {}
}
