package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	"github.com/pquerna/ffjson/ffjson"
)

type ServerConfig struct {
	Port string
}

type HttpStatus struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (httpStatus *HttpStatus) OK(message string) {
	httpStatus.Success = true
	httpStatus.Message = message
}

func (httpStatus *HttpStatus) ERROR(message string) {
	httpStatus.Success = false
	httpStatus.Message = message
}

func JSON(w http.ResponseWriter, structure interface{}) {
	b, err := ffjson.Marshal(structure)
	if err != nil {
		fmt.Fprintf(w, "Failed to parse json: %s", err)
		return
	}
	fmt.Fprintf(w, string(b))
}

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
