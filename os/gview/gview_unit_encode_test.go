// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package gview_test

import (
	"testing"

	"github.com/dxasu/gf/debug/gdebug"
	"github.com/dxasu/gf/frame/g"
	"github.com/dxasu/gf/os/gfile"
	"github.com/dxasu/gf/os/gview"
	"github.com/dxasu/gf/test/gtest"
)

func Test_Encode_Parse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetPath(gdebug.TestDataPath("tpl"))
		v.SetAutoEncode(true)
		result, err := v.Parse("encode.tpl", g.Map{
			"title": "<b>my title</b>",
		})
		t.Assert(err, nil)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}

func Test_Encode_ParseContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		tplContent := gfile.GetContents(gdebug.TestDataPath("tpl", "encode.tpl"))
		v.SetAutoEncode(true)
		result, err := v.ParseContent(tplContent, g.Map{
			"title": "<b>my title</b>",
		})
		t.Assert(err, nil)
		t.Assert(result, "<div>&lt;b&gt;my title&lt;/b&gt;</div>")
	})
}
