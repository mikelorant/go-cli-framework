package config

import (
	"fmt"
	"os"

	"github.com/jeremywohl/flatten/v2"
	"gopkg.in/yaml.v3"
)

func File(filename string) Loader {
	return LoaderFunc(func(cs ConfigSet) error {
		data, err := os.ReadFile(filename)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("unable to read file: %w", err)
		}

		var decodedData map[string]any
		if err := yaml.Unmarshal(data, &decodedData); err != nil {
			return fmt.Errorf("unable to unmarshal yaml: %w", err)
		}

		flatDecodedData, err := flatten.Flatten(decodedData, "", flatten.DotStyle)
		if err != nil {
			fmt.Errorf("unable to flatten config file: %w", err)
		}

		for k, v := range flatDecodedData {
			if _, ok := cs[k]; !ok {
				cs[k] = &ConfigValue{}
			}
			cs[k].value = v
			cs[k].origin = ConfigFile
		}

		return nil
	})
}
