package main

import (
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/net/ghttp"
)

func main() {
	g.I18n().SetPath("/Users/john/Workspace/Go/GOPATH/src/github.com/dxasu/gf/.example/i18n/gi18n/i18n")
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteTplContent(`{#hello}{#world}!`, g.Map{
			"I18nLanguage": r.Get("lang", "zh-CN"),
		})
	})
	s.SetPort(8199)
	s.Run()
}
