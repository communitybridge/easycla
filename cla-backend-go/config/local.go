// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func loadLocalConfig(configFilePath string) (Config, error) {
	configData, err := ioutil.ReadFile(filepath.Clean(configFilePath))
	if err != nil {
		return Config{}, err
	}

	localConfig := Config{}
	err = json.Unmarshal(configData, &localConfig)
	if err != nil {
		return Config{}, err
	}

	return localConfig, nil
}
