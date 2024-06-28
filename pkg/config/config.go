package config

import (
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func NewConfig(configPath string) *Bootstrap {
	k := koanf.New(".")
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		panic(err)
	}

	c := Bootstrap{}
	if err := k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		panic(err)
	}

	return &c
}
