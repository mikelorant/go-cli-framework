package config

func Defaults(bindings Bindings) Loader {
	return LoaderFunc(func(cs ConfigSet) error {
		for _, v := range bindings {
			if v.flag.DefValue != "" {
				if _, ok := cs[v.key]; !ok {
					cs[v.key] = &ConfigValue{}
				}
				cs[v.key].value = v.flag.DefValue
				cs[v.key].origin = ConfigDefault
			}
		}

		return nil
	})
}
