package main

import (
	"github.com/dxasu/gf/util/gutil"
)

func Test(s *interface{}) {
	//debug.PrintStack()
	gutil.PrintStack()
}

func main() {
	Test(nil)
}
