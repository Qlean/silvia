package main

import(
	"log"
	"net/http"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port      string
	PgConnect string `yaml:"pg_connect"`
	RingSize  int    `yaml:"ring_size"`
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

func eventFail(w http.ResponseWriter, localCounter int, errMessage string, err error, body string, nginxLog *NginxLog, event *Event) {
	message := errMessage + ": " + err.Error()
	http.Error(w, message, http.StatusBadRequest)
	log.Println(localCounter, message)
	failRing.Add(body, nginxLog, event, message)
	failCounter++
}

