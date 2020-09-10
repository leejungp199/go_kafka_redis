package module

import (
	"./config"
	"gopkg.in/ini.v1"
)

var (
	cfg *config.Config
)

func initConfig(configPath string) *config.Config {
	c := new(config.Config)
	err := ini.MapTo(c, configPath)
	if err != nil {
		panic(err)
	}
	return c
}
