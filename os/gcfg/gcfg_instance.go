// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package gcfg

import (
	"fmt"

	"github.com/dxasu/gf/container/gmap"
)

const (
	// Default group name for instance usage.
	DefaultName = "config"
)

var (
	// Instances map containing configuration instances.
	instances = gmap.NewStrAnyMap(true)
)

// Instance returns an instance of Config with default settings.
// The parameter `name` is the name for the instance. But very note that, if the file "name.toml"
// exists in the configuration directory, it then sets it as the default configuration file. The
// toml file type is the default configuration file type.
func Instance(name ...string) *Config {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		// If it's not using default configuration or its configuration file is not available,
		// it searches the possible configuration file according to the name and all supported
		// file types.
		if key != DefaultName || !c.Available() {
			for _, fileType := range supportedFileTypes {
				if file := fmt.Sprintf(`%s.%s`, key, fileType); c.Available(file) {
					c.SetFileName(file)
					break
				}
			}
		}
		return c
	}).(*Config)
}
