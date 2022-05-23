package config

import (
	"fmt"
	"os"
	"strings"
)

func Env(bindings Bindings, prefix string) Loader {
	return LoaderFunc(func(cs ConfigSet) error {
		for _, v := range bindings {
			envDotted := strings.ReplaceAll(v.key, ".", "_")
			envUpper := strings.ToUpper(envDotted)
			envPrefixed := fmt.Sprintf("%v_%v", prefix, envUpper)

			if env, ok := os.LookupEnv(envPrefixed); ok {
				if _, ok := cs[v.key]; !ok {
					cs[v.key] = &ConfigValue{}
				}
				cs[v.key].value = env
				cs[v.key].origin = ConfigEnv
			}
		}

		return nil
	})
}
