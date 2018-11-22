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
	"path"

	"github.com/alecthomas/units"
	"github.com/spf13/viper"
)

func readConfigFile(customFile ...string) {
	if len(customFile) > 0 {
		var filename = path.Base(customFile[0])
		var fileDir = path.Dir(customFile[0])
		viper.SetConfigFile(filename)
		viper.AddConfigPath(fileDir)
		err := viper.ReadInConfig() // Find and read the config file
		// If there is no configuration file, that's an error
		if err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigName("cn")         // name of config file (without extension)
		viper.AddConfigPath("/etc/ceph/") // path to look for the config file in
		viper.AddConfigPath("$HOME/.cn/") // call multiple times to add many search paths
		viper.AddConfigPath(".")          // optionally look for config in the working directory
		viper.ReadInConfig()              // Find and read the config file, we don't really care if no config file is found
	}
	setDefaultConfig()
}

// Set the default values for defined types
// If the configuration file is missing, this section will generated the mandatory elements
func setDefaultConfig() {
	viper.SetDefault("default.MemorySize", "512MB")
}

func getValueFromConfig(name string, cluster ...string) string {
	var keyname = "default" + "." + name
	// If a clustername is given let's override the keyname with it
	if len(cluster) == 1 {
		keyname = cluster[0] + "." + name
	}
	return viper.GetString(keyname)
}

func getMemorySize(containerName ...string) int64 {
	var bytes units.Base2Bytes
	var err error
	if len(containerName) > 0 {
		bytes, err = units.ParseBase2Bytes(getValueFromConfig("MemorySize", containerName[0]))
	} else {
		bytes, err = units.ParseBase2Bytes(getValueFromConfig("MemorySize"))
	}
	if err != nil {
		panic(err)
	}
	return int64(bytes)
}
