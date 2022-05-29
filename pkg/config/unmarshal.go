package config

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/nqd/flat"
)

func (c *Config) Unmarshal(out any) error {
	return Unmarshal(c.AllSettings(), out)
}

func (c *Config) UnmarshalWithKey(key string, out any) error {
	return UnmarshalWithKey(c.AllSettings(), key, out)
}

func Unmarshal(in map[string]any, out any) error {
	nested, err := flat.Unflatten(in, nil)
	if err != nil {
		fmt.Errorf("unable to unflatten map: %w", err)
	}

	if err := mapstructure.Decode(nested, out); err != nil {
		return fmt.Errorf("unable to decode map: %w", err)
	}

	return nil
}

func UnmarshalWithKey(in map[string]any, key string, out any) error {
	sub := getKey(in, key)
	if err := Unmarshal(sub, out); err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
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
