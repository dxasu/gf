package testdata

import "github.com/dxasu/gf/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
