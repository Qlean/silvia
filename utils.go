package main

import(
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port      string
	Pgconnect string
}

var serverConfig *ServerConfig

func LoadConfig(filename string) (server_config *ServerConfig, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(buf, &server_config)
	if err != nil {
		return
	}

	return
}
