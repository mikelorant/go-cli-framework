package config

func Defaults(bindings Bindings) Loader {
	return LoaderFunc(func(cs ConfigSet) error {
		for _, v := range bindings {
			if !v.flag.Changed {
				if _, ok := cs[v.key]; !ok {
					cs[v.key] = &ConfigValue{}
				}
				cs[v.key].value = v.flag.Value
				cs[v.key].origin = ConfigDefault
			}
		}

		return nil
	})
}
