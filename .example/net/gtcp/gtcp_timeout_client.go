package main

import (
	"time"

	"github.com/dxasu/gf/net/gtcp"
	"github.com/dxasu/gf/os/glog"
	"github.com/dxasu/gf/os/gtime"
)

func main() {
	conn, err := gtcp.NewConn("127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := conn.Send([]byte(gtime.Now().String())); err != nil {
		glog.Error(err)
	}

	time.Sleep(time.Minute)
}
