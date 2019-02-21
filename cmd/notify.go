/*
 * Ceph Nano (C) 2019 Red Hat, Inc.
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
	"io/ioutil"
	"log"
	"time"
)

var (
	timeLayout              = time.RFC1123
	lastUpdateCheckFilePath = makeCephNanoPath("last_update_check")
)

func checkUpdateNotification() {
	if !shouldCheckURLVersion(lastUpdateCheckFilePath) {
		return
	}
	updateCheckNano(nil, nil)
	writeTimeToFile(lastUpdateCheckFilePath, time.Now().UTC())
}

func shouldCheckURLVersion(filePath string) bool {
	if !getBoolFromConfig("update", "config", "want_update_notification") {
		return false
	}
	lastUpdateTime := getTimeFromFileIfExists(filePath)
	return time.Since(lastUpdateTime).Hours() >= getFloat64FromConfig("update", "config", "reminder_wait_period_in_hours")
}

func writeTimeToFile(path string, inputTime time.Time) {
	err := ioutil.WriteFile(path, []byte(inputTime.Format(timeLayout)), 0644)
	if err != nil {
		log.Fatalf("Error writing current update time to file: %v", err)
	}
}

func getTimeFromFileIfExists(path string) time.Time {
	lastUpdateCheckTime, err := ioutil.ReadFile(path)
	if err != nil {
		return time.Time{}
	}
	timeInFile, err := time.Parse(timeLayout, string(lastUpdateCheckTime))
	if err != nil {
		return time.Time{}
	}
	return timeInFile
}
