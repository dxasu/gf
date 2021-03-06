package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
		body1 := r.GetBody()
		body2, _ := ioutil.ReadAll(r.Body)
		fmt.Println(body1)
		fmt.Println(body2)
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}
