// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package ghttp_test

import (
	"github.com/dxasu/gf/container/garray"
	"github.com/dxasu/gf/os/genv"
)

var (
	ports = garray.NewIntArray(true)
)

func init() {
	genv.Set("UNDER_TEST", "1")
	for i := 7000; i <= 8000; i++ {
		ports.Append(i)
	}
}
