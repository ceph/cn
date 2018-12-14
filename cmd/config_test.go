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

func TestDefaultConfig(t *testing.T) {
	// Testing the builtin configuration
	assert.Equal(t, "512MB", getMemorySize("default"))
	assert.Equal(t, "4GB", getMemorySize("huge"))
	assert.Equal(t, int64(1), getCPUCount("default"))
	assert.Equal(t, int64(2), getCPUCount("huge"))
	assert.Equal(t, DEFAULTWORKDIRECTORY, getWorkDirectory("default"))
	assert.Equal(t, true, isEntryExists(FLAVORS, "default.use_default"))
	assert.Equal(t, false, isEntryExists(FLAVORS, "default.nawak"))
	assert.Equal(t, false, getPrivileged("default"))
	setPrivileged("default", true)
	assert.Equal(t, true, getPrivileged("default"))
	setPrivileged("default", false)
	assert.Equal(t, "", getSize("default"))

	defaultImageName := imageName
	// Without any configuration file, the default should be satisfied
	assert.Equal(t, DEFAULTIMAGE, getImageName())

	// Without any configuration file, any -i argument should be preserved
	imageName = "nawak"
	assert.Equal(t, "nawak", getImageName())

	// The default builtin should be kept too
	imageName = "mimic"
	assert.Equal(t, LATESTIMAGE+"mimic", getImageName())

	imageName = defaultImageName
	assert.Equal(t, "", getUnderlyingStorage("default"))
}

func TestReadConfigFile(t *testing.T) {
	assert.Equal(t, configFile, readConfigFile(configFile))
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// STARTING HERE ALL TESTS ARE RUN AGAINST A CONFIGURATION FILE
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

func TestTitle(t *testing.T) {
	assert.Equal(t, "Ceph Nano test configuration file", viper.Get("title"))
}

func TestMemorySize(t *testing.T) {
	assert.Equal(t, "512MB", getMemorySize("test_nano_default"))
	assert.Equal(t, int64(536870912), getMemorySizeInBytes("test_nano_default"))
	assert.Equal(t, "1GB", getMemorySize("test_nano_no_default"))
	assert.Equal(t, int64(1073741824), getMemorySizeInBytes("test_nano_no_default"))
}

func TestUseDefault(t *testing.T) {
	assert.Equal(t, false, useDefault(FLAVORS, "test_nano_no_default"))
	assert.Equal(t, true, useDefault(FLAVORS, "test_nano_default"))
}

func TestCephConf(t *testing.T) {
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
	assert.Equal(t, int64(1), getCPUCount("test_nano_default"))
	assert.Equal(t, int64(2), getCPUCount("test_nano_no_default"))
}

func TestPrivileged(t *testing.T) {
	assert.Equal(t, true, getPrivileged("test_nano_no_default"))
}

func TestImageName(t *testing.T) {
	defaultImageName := imageName
	// Let's ensure the basic reading of the configuration file works
	assert.Equal(t, "ceph/daemon:latest-real1", getImageNameFromConfig("real1"))

	// If a -i is passed with a configuration file, let's report the image_name from the configuration file
	imageName = "complex"
	assert.Equal(t, "this.url.is.complex/cool/for-a-test", getImageName())
	imageName = defaultImageName
}

func TestUnderlyingStorage(t *testing.T) {
	defaultDataOsd := dataOsd
	// Ensure the values are properly read from the configuration file
	assert.Equal(t, "/dev/sdb1", getUnderlyingStorage("test_nano_no_default"))

	// Ensure that enforcing a flags on the CLI is taking over the configuration
	dataOsd = "/dev/nawak"
	assert.Equal(t, "/dev/nawak", getUnderlyingStorage("test_nano_no_default"))
	dataOsd = defaultDataOsd
}

func TestSize(t *testing.T) {
	defaultSizeBluestoreBlock := sizeBluestoreBlock
	// Ensure the values are properly read from the configuration file
	assert.Equal(t, "20GB", getSize("test_nano_no_default"))

	// Ensure that enforcing a flags on the CLI is taking over the configuration
	sizeBluestoreBlock = "1M"
	assert.Equal(t, "1M", getSize("test_nano_no_default"))
	sizeBluestoreBlock = defaultSizeBluestoreBlock
}

func TestWorkDirectory(t *testing.T) {
	defaultWorkingDirectory := workingDirectory
	// Ensure the values are properly read from the configuration file
	assert.Equal(t, "/tmp/nano/", getWorkDirectory("test_nano_no_default"))

	// Ensure that enforcing a flags on the CLI is taking over the configuration
	workingDirectory = "/tmp/nawak"
	assert.Equal(t, "/tmp/nawak", getWorkDirectory("test_nano_no_default"))
	workingDirectory = defaultWorkingDirectory
}

func TestMerge(t *testing.T) {
	assert.Equal(t, false, isParameterExist(FLAVORS, "test_nano_no_default", "new_param"))
	assert.Equal(t, true, isParameterExist(FLAVORS, "test_nano_default", "new_param"))
}
