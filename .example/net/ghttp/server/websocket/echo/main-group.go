package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
	"github.com/dxasu/gf/os/gfile"
	"github.com/dxasu/gf/os/glog"
)

func ws(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		return
	}
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {
	s := g.Server()
	s.Group("").Bind([]ghttp.GroupItem{
		{"ALL", "/ws", ws},
	})
	s.SetAccessLogEnabled(true)
	s.SetServerRoot(gfile.MainPkgPath())
	s.SetPort(8199)
	s.Run()
}
