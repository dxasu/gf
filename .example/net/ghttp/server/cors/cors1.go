package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
	"github.com/dxasu/gf/os/glog"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Order(r *ghttp.Request) {
	glog.Println("order")
	r.Response.Write("GET")
}

func main() {
	s := g.Server()
	s.Group("/api.v1", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		group.GET("/order", Order)
	})
	s.SetPort(8199)
	s.Run()
}
