package main

import (
	"errors"

	"github.com/dxasu/gf/errors/gerror"
	"github.com/dxasu/gf/os/glog"
)

func MakeError() error {
	return errors.New("connection closed with normal error")
}

func MakeGError() error {
	return gerror.New("connection closed with gerror")
}

func TestGError() {
	err1 := MakeError()
	err2 := MakeGError()
	glog.Error(err1)
	glog.Error(err2)
}

func main() {
	TestGError()
}
