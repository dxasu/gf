package main

import (
	"fmt"
	"path/filepath"

	"github.com/dxasu/gf/frame/g"

	"github.com/dxasu/gf/net/ghttp"
	"github.com/dxasu/gf/os/glog"
)

func SendXmlFile(gameId int, areaName string, filePath string) error {
	path := filepath.FromSlash(filePath)
	fmt.Println(path)
	data := g.Map{
		"gameName":    gameId,
		"area":        areaName,
		"file":        "@file:" + path,
		"contentType": "json",
	}
	if r, err := ghttp.Post("http://127.0.0.1:8199/upload", data); err != nil {
		panic(err)
	} else {
		defer r.Close()
		fmt.Println("ok")
	}
	return nil
}

func main() {
	SendXmlFile(1, "xxx", "/Users/john/Workspace/Go/GOPATH/src/github.com/dxasu/gf/.example/net/ghttp/server/session.go")
	return
	path := "/home/john/Workspace/Go/github.com/dxasu/gf/version.go"
	r, e := ghttp.Post("http://127.0.0.1:8199/upload", "upload-file=@file:"+path)
	if e != nil {
		glog.Error(e)
	} else {
		fmt.Println(string(r.ReadAll()))
		r.Close()
	}
}
