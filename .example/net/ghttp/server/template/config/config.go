package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTplContent(`${.name}`, g.Map{
			"name": "john",
		})
	})
	s.SetPort(8199)
	s.Run()
}
