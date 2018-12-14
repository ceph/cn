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
	"log"
	"path"
	"strings"

	"github.com/spf13/viper"
)

//FLAVORS is a constant to represent the [flavors] group
const FLAVORS = "flavors"

//IMAGES is a constant to represent the [images] group
const IMAGES = "images"

// DEFAULTIMAGE is the default image name to be used
const DEFAULTIMAGE = "ceph/daemon"

// LATESTIMAGE is the prefix for the latest ceph images
const LATESTIMAGE = DEFAULTIMAGE + ":latest-"

// DEFAULTWORKDIRECTORY is the default work directory
const DEFAULTWORKDIRECTORY = "/usr/share/ceph-nano"

func readConfigFile(customFile ...string) string {
	setDefaultConfig()

	// A custom configuration file got passed
	// Let's handle it directly
	if len(customFile) > 0 {
		var filename = path.Base(customFile[0])
		var fileDir = path.Dir(customFile[0])
		viper.SetConfigFile(filename)
		viper.AddConfigPath(fileDir)
		err := viper.ReadInConfig()
		// Find and read the config file
		// If there is no configuration file, that's an error
		if err != nil {
			log.Fatal(err)
		}

		// Let's import all the default value into flavors
		mergeFlavorsWithDefault()

		// We return the configuration file found
		return viper.ConfigFileUsed()
	}

	// Let's search for a configuration file
	viper.SetConfigName("cn")         // name of config file (without extension)
	viper.AddConfigPath("/etc/cn/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.cn/") // call multiple times to add many search paths
	viper.AddConfigPath(".")          // optionally look for config in the working directory
	// Let's try to read a configuration file
	if viper.ReadInConfig() == nil {
		// Let's import all the default value into flavors
		mergeFlavorsWithDefault()
		// We return the configuration file found
		return viper.ConfigFileUsed()
	}

	// No configuration file is used aka compatiblity mode
	return ""
}

// Set the default values for defined types
// If the configuration file is missing, this section will generated the mandatory elements
func setDefaultConfig() {
	// Handling the built-in flavor
	viper.SetDefault(FLAVORS+".default.use_default", true) // All containers inherit from default
	viper.SetDefault(FLAVORS+".default.memory_size", "512MB")
	viper.SetDefault(FLAVORS+".default.cpu_count", 1)
	viper.SetDefault(FLAVORS+".default.privileged", false)
	viper.SetDefault(FLAVORS+".default.data", "")
	viper.SetDefault(FLAVORS+".default.size", "")
	viper.SetDefault(FLAVORS+".default.work_directory", DEFAULTWORKDIRECTORY)
	viper.SetDefault(FLAVORS+".medium.memory_size", "768MB")
	viper.SetDefault(FLAVORS+".large.memory_size", "1GB")
	viper.SetDefault(FLAVORS+".huge.memory_size", "4GB")
	viper.SetDefault(FLAVORS+".huge.cpu_count", 2)

	// Handling the built-in image aliases
	viper.SetDefault(IMAGES+".default.use_default", true) // All containers inherit from default
	viper.SetDefault(IMAGES+".default.image_name", DEFAULTIMAGE)
	// Setting up the aliases to be reported in 'image show-aliases' command
	viper.SetDefault(IMAGES+".mimic.image_name", LATESTIMAGE+"mimic")
	viper.SetDefault(IMAGES+".luminous.image_name", LATESTIMAGE+"luminous")
	viper.SetDefault(IMAGES+".redhat.image_name", "registry.access.redhat.com/rhceph/rhceph-3-rhel7")
}

func getStringFromConfig(group string, item string, name string) string {
	var value = ""
	// If we are requested to get the status of use_default, we cannot call useDefault ;)
	itemValue := viper.GetString(group + "." + item + "." + name)
	if len(itemValue) > 0 {
		value = itemValue
	}
	return value
}

func getInt64FromConfig(group string, item string, name string) int64 {
	var value int64
	var foundValue = false

	// We need to ensure the key exists unless that could populate a 0 value
	if isParameterExist(group, item, name) {
		value = viper.GetInt64(group + "." + item + "." + name)
		foundValue = true
	}

	if !foundValue {
		log.Fatal(name + " int64 value in " + item + "doesn't exist")
	}
	return value
}

func getBoolFromConfig(group string, item string, name string) bool {
	// We need to ensure the key exist unless that could populate a wrong value
	if isParameterExist(group, item, name) {
		return viper.GetBool(group + "." + item + "." + name)
	}

	// If we are reaching this point, let's check if we are in the chicken/egg case trigger by mergeFlavorsWithDefault
	// As mergeFlavorsWithDefault checks if use_default is while its not yet populated
	// In such case, we should read the default first
	if name == "use_default" {
		if isParameterExist(group, "default", name) {
			return viper.GetBool(group + ".default." + name)
		}
	}

	// If we reach that point, that means this bool doesn't exist
	log.Fatal(name + " bool value in " + item + " doesn't exist")
	// We cannot reach this point
	return false
}

func useDefault(group string, item string) bool {
	return getBoolFromConfig(group, item, "use_default")
}

func getStringMapFromConfig(group string, item string, name string) map[string]interface{} {
	var defaultConfig = make(map[string]interface{})
	if useDefault(group, item) {
		defaultConfig = viper.GetStringMap(group + ".default" + "." + name)
	}
	itemValues := viper.GetStringMap(group + "." + item + "." + name)
	if len(itemValues) > 0 {
		for key, value := range itemValues {
			defaultConfig[key] = value
		}
	}
	return defaultConfig
}

func isEntryExists(group string, item string) bool {
	return viper.IsSet(group + "." + item)
}

// Return items from a given group
func getItemsFromGroup(group string) map[string]interface{} {
	return viper.AllSettings()[group].(map[string]interface{})
}

// Does this parameter exist in the configuration
func isParameterExist(group string, item string, parameter string) bool {
	return viper.IsSet(group + "." + item + "." + parameter)
}

// Afunction to list the default parameters as they are not seen
func getDefaultParameters() map[string]interface{} {
	returnValue := make(map[string]interface{})
	// For each keys in the configuration
	for _, param := range viper.AllKeys() {
		// If there is a default entry
		if strings.HasPrefix(param, FLAVORS+".default.") {
			// extract the parameter name
			parameter := strings.SplitAfter(param, FLAVORS+".default.")[1]
			// Let's return the association parameter/value
			returnValue[parameter] = viper.Get(param)
		}
	}
	return returnValue
}

// Considering the use_default value, let's merge the default values in other flavors
func mergeFlavorsWithDefault() {
	// For each flavor
	for flavor := range getItemsFromGroup(FLAVORS) {

		// Let's skip the default flavor
		if flavor == "default" {
			continue
		}

		// Nothing to do if the flavor set the use_default=false
		if !useDefault(FLAVORS, flavor) {
			continue
		}
		// For every default parameter
		for defaultParameter, defaultValue := range getDefaultParameters() {
			// If the flavor doesn't define it
			if viper.Get(FLAVORS+"."+flavor+"."+defaultParameter) == nil {
				// Let's copy the default value in this flavor
				viper.SetDefault(FLAVORS+"."+flavor+"."+defaultParameter, defaultValue)
			}
		}
	}
}
