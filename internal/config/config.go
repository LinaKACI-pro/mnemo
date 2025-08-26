package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const DELIMITER = "."

var k = koanf.New(DELIMITER)

type Config struct {
	DbPath   string `koanf:"db_path"`
	DbDriver string `koanf:"db_driver"`
}

// Load configuration with the following precedence:
// 1. Defaults (hardcoded)
// 2. Optional config.yaml file
// 3. Environment variables (prefix: MNEMO_)
func Load() (Config, error) {
	// 1. Defaults
	k.Load(confmap.Provider(map[string]interface{}{
		"db_path":   "./mnemo.db",
		"db_driver": "sqlite",
	}, DELIMITER), nil)

	// 2. Optional config.yaml file
	_ = k.Load(file.Provider("config.yaml"), yaml.Parser())

	// 3. Environment variables (MNEMO_DB_PATH)
	k.Load(env.Provider(".", env.Opt{
		Prefix: "MNEMO_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ToLower(strings.TrimPrefix(k, "MNEMO_"))

			// If there is a space in the value, split the value into a slice by the space.
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}
			return k, v
		},
		EnvironFunc: func() []string {
			return slices.DeleteFunc(os.Environ(), func(s string) bool {
				return strings.HasPrefix(s, "MNEMO_TIME")
			})
		},
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
