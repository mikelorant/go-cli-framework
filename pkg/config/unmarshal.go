package config

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func (c *Config) Unmarshal(out any) error {
	return Unmarshal(c.AllSettings(), out)
}

func (c *Config) UnmarshalWithKey(key string, out any) error {
	return UnmarshalWithKey(c.AllSettings(), key, out)
}

func Unmarshal(in any, out any) error {
	if err := mapstructure.Decode(in, out); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}

func UnmarshalWithKey(in map[string]any, key string, out any) error {
	if err := mapstructure.Decode(getKey(in, key), out); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	return nil
}

func getKey(in map[string]any, key string) map[string]any {
	out := make(map[string]any)

	prefix := fmt.Sprintf("%v%v", key, ".")

	for k, v := range in {
		if strings.HasPrefix(k, key) {
			out[strings.TrimPrefix(k, prefix)] = v
		}
	}

	return out
}
