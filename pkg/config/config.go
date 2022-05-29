package config

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
)

type Config struct {
	bindings  Bindings
	configSet ConfigSet
	options   *Options
}

type ConfigAdapter func(*Config)

type Bindings []*Binding

type Binding struct {
	key  string
	flag *pflag.Flag
}

type ConfigSet map[string]*ConfigValue

type ConfigValue struct {
	value  any
	origin Origin
}

type Origin int

type Options struct {
	flags     *pflag.FlagSet
	envPrefix string
	filename  string
}

type Loader interface {
	Process(cs ConfigSet) error
}

type LoaderFunc func(cs ConfigSet) error

const (
	ConfigUnset Origin = iota
	ConfigDefault
	ConfigFile
	ConfigEnv
	ConfigFlag
)

func New(ctx context.Context) *Config {
	return &Config{
		options:   &Options{},
		bindings:  Bindings{},
		configSet: ConfigSet{},
	}
}

func WithFlags(flags *pflag.FlagSet) ConfigAdapter {
	return func(c *Config) {
		c.options.flags = flags
	}
}

func WithEnvPrefix(prefix string) ConfigAdapter {
	return func(c *Config) {
		c.options.envPrefix = prefix
	}
}

func WithFilename(filename string) ConfigAdapter {
	return func(c *Config) {
		c.options.filename = filename
	}
}

func (fn LoaderFunc) Process(cs ConfigSet) error {
	return fn(cs)
}

func (c *Config) Load(opts ...ConfigAdapter) error {
	for _, o := range opts {
		o(c)
	}

	if err := c.bind(); err != nil {
		return fmt.Errorf("unable to bind flags: %w", err)
	}

	Defaults(c.bindings).Process(c.configSet)
	File(c.options.filename).Process(c.configSet)
	Env(c.bindings, c.options.envPrefix).Process(c.configSet)
	Flag(c.bindings).Process(c.configSet)

	return nil
}

func (c *Config) bind() error {
	const bindAnnotation = "bindWithKey"

	c.options.flags.VisitAll(func(flag *pflag.Flag) {
		if _, ok := flag.Annotations[bindAnnotation]; !ok {
			return
		}

		b := &Binding{
			key:  flag.Annotations[bindAnnotation][0],
			flag: flag,
		}
		c.bindings = append(c.bindings, b)
	})

	return nil
}

func (c *Config) AllSettings() map[string]any {
	cfg := make(map[string]any)

	for k, v := range c.configSet {
		cfg[k] = v.value
	}

	return cfg
}

func (c *Config) AllSettingsChanged() map[string]any {
	cfg := make(map[string]any)

	for k, v := range c.configSet {
		if v.origin > ConfigDefault {
			cfg[k] = v.value
		}
	}

	return cfg
}
