/*
 * Ceph Nano (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var configFile = "cn-test.toml"

func TestTitle(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, "Ceph Nano test configuration file", viper.Get("title"))
}

func TestMemorySize(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, int64(536870912), getMemorySize())
	assert.Equal(t, int64(1073741824), getMemorySize("test_nano_no_default"))
}

func TestUseDefault(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, false, useDefault("test_nano_no_default"))
	assert.Equal(t, true, useDefault("test_nano_default"))
}
