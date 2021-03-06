// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package gcfg

import (
	"github.com/dxasu/gf/container/gmap"
	"github.com/dxasu/gf/internal/intlog"
)

var (
	// Customized configuration content.
	configs = gmap.NewStrStrMap(true)
)

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func SetContent(content string, file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	instances.LockFunc(func(m map[string]interface{}) {
		if configs.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
		}
		configs.Set(name, content)
	})
}

// GetContent returns customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func GetContent(file ...string) string {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	return configs.Get(name)
}

// RemoveContent removes the global configuration with specified `file`.
// If `name` is not passed, it removes configuration of the default group name.
func RemoveContent(file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached `name`.
	instances.LockFunc(func(m map[string]interface{}) {
		if configs.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
			configs.Remove(name)
		}
	})

	intlog.Printf(`RemoveContent: %s`, name)
}

// ClearContent removes all global configuration contents.
func ClearContent() {
	configs.Clear()
	// Clear cache for all instances.
	instances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			v.(*Config).jsonMap.Clear()
		}
	})

	intlog.Print(`RemoveConfig`)
}
