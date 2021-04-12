package main

import (
	"github.com/dxasu/gf/os/gview"
	"github.com/dxasu/gf/util/gutil"
)

func main() {
	gutil.Dump(gview.ParseContent(`{{"<div>测试</div>去掉HTML标签<script>var v=1;</script>"|text}}`, nil))
}
