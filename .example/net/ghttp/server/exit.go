package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
	"github.com/dxasu/gf/os/glog"
)

func main() {
	p := "/"
	s := g.Server()
	s.BindHandler(p, func(r *ghttp.Request) {
		r.Response.Writeln("start")
		r.Exit()
		r.Response.Writeln("end")
	})
	s.BindHookHandlerByMap(p, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			glog.To(r.Response.Writer).Println("BeforeServe")
		},
		ghttp.HookAfterServe: func(r *ghttp.Request) {
			glog.To(r.Response.Writer).Println("AfterServe")
		},
	})
	s.SetPort(8199)
	s.Run()
}
