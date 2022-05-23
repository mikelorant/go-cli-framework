package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jeremywohl/flatten/v2"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Config struct {
	options *Options
	binds   []*Bind
	config  map[string]*ConfigValue
}

type ConfigAdapter func(*Config)

type Loader interface {
	Process(dst any) error
}

type LoaderFunc func(dst any) error

type Options struct {
	flags     *pflag.FlagSet
	envPrefix string
	filename  string
}

type ConfigValue struct {
	value  any
	origin Origin
}

type Bind struct {
	key  string
	flag *pflag.Flag
}

type Origin int

const (
	ConfigUnset Origin = iota
	ConfigDefault
	ConfigFile
	ConfigEnv
	ConfigFlag
)

func New(ctx context.Context) *Config {
	return &Config{
		options: &Options{},
		binds:   []*Bind{},
		config:  make(map[string]*ConfigValue),
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

func (c *Config) AllSettings() map[string]any {
	cfg := make(map[string]any)

	for k, v := range c.config {
		cfg[k] = v.value
	}

	return cfg
}

func (c *Config) AllSettingsChanged() map[string]any {
	cfg := make(map[string]any)

	for k, v := range c.config {
		if v.origin > ConfigDefault {
			cfg[k] = v.value
		}
	}

	return cfg
}

func (c *Config) Load(opts ...ConfigAdapter) error {
	for _, o := range opts {
		o(c)
	}

	if err := c.bind(); err != nil {
		return fmt.Errorf("unable to bind flags: %w", err)
	}

	if err := c.loadDefault(); err != nil {
		return fmt.Errorf("unable to load defaults: %w", err)
	}

	if err := c.loadFile(c.options.filename); err != nil {
		return fmt.Errorf("unable to load file: %w", err)
	}

	if err := c.loadEnv(c.options.envPrefix); err != nil {
		return fmt.Errorf("unable to load env: %w", err)
	}

	if err := c.loadFlag(); err != nil {
		return fmt.Errorf("unable to load flags: %w", err)
	}

	return nil
}

func (c *Config) loadDefault() error {
	for _, v := range c.binds {
		if v.flag.DefValue != "" {
			if _, ok := c.config[v.key]; !ok {
				c.config[v.key] = &ConfigValue{}
			}
			c.config[v.key].value = v.flag.DefValue
			c.config[v.key].origin = ConfigDefault
		}
	}

	return nil
}

func (c *Config) loadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	var cfg map[string]any
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	cfgflat, err := flatten.Flatten(cfg, "", flatten.DotStyle)
	if err != nil {
		fmt.Errorf("unable to flatten config file: %w", err)
	}

	for k, v := range cfgflat {
		if _, ok := c.config[k]; !ok {
			c.config[k] = &ConfigValue{}
		}
		c.config[k].value = v
		c.config[k].origin = ConfigFile
	}

	return nil
}

func (c *Config) loadEnv(prefix string) error {
	for _, v := range c.binds {
		envDotted := strings.ReplaceAll(v.key, ".", "_")
		envUpper := strings.ToUpper(envDotted)
		envPrefixed := fmt.Sprintf("%v_%v", prefix, envUpper)

		if env, ok := os.LookupEnv(envPrefixed); ok {
			if _, ok := c.config[v.key]; !ok {
				c.config[v.key] = &ConfigValue{}
			}
			c.config[v.key].value = env
			c.config[v.key].origin = ConfigEnv
		}
	}

	return nil
}

func (c *Config) loadFlag() error {
	for _, v := range c.binds {
		if v.flag.Changed {
			if _, ok := c.config[v.key]; !ok {
				c.config[v.key] = &ConfigValue{}
			}
			c.config[v.key].value = v.flag.Value
			c.config[v.key].origin = ConfigFlag
		}
	}

	return nil
}

func (c *Config) bind() error {
	const bindAnnotation = "bindWithKey"

	c.options.flags.VisitAll(func(flag *pflag.Flag) {
		if _, ok := flag.Annotations[bindAnnotation]; !ok {
			return
		}

		b := &Bind{
			key:  flag.Annotations[bindAnnotation][0],
			flag: flag,
		}
		c.binds = append(c.binds, b)
	})

	return nil
}
