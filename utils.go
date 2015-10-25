package main

import(
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port       string
	PgConnect  string `yaml:"pg_connect"`
	RingSize   int    `yaml:"ring_size"`
	RabbitAddr string `yaml:"rabbit_addr"`
	RabbitPort string `yaml:"rabbit_port"`
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

