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
	assert.Equal(t, int64(536870912), getMemorySize("test_nano_default"))
	assert.Equal(t, int64(1073741824), getMemorySize("test_nano_no_default"))
}

func TestUseDefault(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, false, useDefault("test_nano_no_default"))
	assert.Equal(t, true, useDefault("test_nano_default"))
}

func TestCephConf(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, map[string]interface{}{"osd_memory_target": int64(3841234556)}, getCephConf("test_nano_no_default"))
	expectedOutput := map[string]interface{}{
		"bluestore_cache_autotune_chunk_size": int64(8388608),
		"osd_max_pg_log_entries":              int64(10),
		"osd_memory_base":                     int64(268435456),
		"osd_memory_cache_min":                int64(33554432),
		"osd_memory_target":                   int64(3841234556),
		"osd_min_pg_log_entries":              int64(10),
		"osd_pg_log_dups_tracked":             int64(10),
		"osd_pg_log_trim_min":                 int64(10),
	}
	assert.Equal(t, expectedOutput, getCephConf("test_nano_default"))
}

func TestCPUCount(t *testing.T) {
	readConfigFile(configFile)
	assert.Equal(t, int64(1), getCPUCount("test_nano_default"))
	assert.Equal(t, int64(2), getCPUCount("test_nano_no_default"))
}
