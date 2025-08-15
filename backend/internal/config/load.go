package config

import (
	"encoding/base64"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"log"
	"log/slog"
	"reflect"
	"strings"
)

const (
	tag       = "config"
	delimiter = "."
	prefix    = "BERMUDIA__"
	separator = "__"
)

func Load() Config {
	k := koanf.New(delimiter)

	{
		err := k.Load(structs.Provider(defaultConfig(), tag), nil)
		if err != nil {
			log.Fatalf("could not load default config: %v\n", err)
		}
	}

	{
		err := k.Load(env.Provider(prefix, delimiter, envCallBack), nil)
		if err != nil {
			slog.Error("could not load env variables for config: %v\n", err)
		}
	}

	var instance Config
	err := k.UnmarshalWithConf("", &instance, koanf.UnmarshalConf{
		Tag: tag,
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				func(f reflect.Type, t reflect.Type, data any) (any, error) {
					if f.Kind() != reflect.String {
						return data, nil
					}
					if t != reflect.TypeFor[[]byte]() {
						return data, nil
					}
					return base64.StdEncoding.DecodeString(data.(string))
				},
			),
		},
	})

	if err != nil {
		log.Fatalf("could not unmarshal config: %v\n", err)
	}

	return instance
}

func envCallBack(s string) string {
	base := strings.ToLower(strings.TrimPrefix(s, prefix))
	return strings.ReplaceAll(base, separator, delimiter)
}
