package main

import (
	"fmt"

	"github.com/dxasu/gf/util/gutil"
)

func main() {
	gutil.TryCatch(func() {
		fmt.Println(1)
		gutil.Throw("error")
		fmt.Println(2)
	}, func(err error) {
		fmt.Println(err)
	})
}
