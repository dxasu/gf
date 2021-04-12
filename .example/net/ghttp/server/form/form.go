package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Dump(r.GetForm("array"))
		r.Response.WriteTpl("form.html")
	})
	s.SetPort(8199)
	s.Run()
}
