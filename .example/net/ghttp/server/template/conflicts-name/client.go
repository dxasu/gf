package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
)

// https://github.com/dxasu/gf/issues/437
func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTpl("client/layout.html")
	})
	s.SetPort(8199)
	s.Run()
}
