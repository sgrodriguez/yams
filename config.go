package main

import "github.com/BurntSushi/toml"

//LoadConfig from toml file
func LoadConfig(tomlPath string) (*Config, error) {
	conf := &Config{}
	_, err := toml.DecodeFile(tomlPath, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Config ...
type Config struct {
	Port  int
	Mocks []Mock
}
