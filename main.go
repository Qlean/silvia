package main

import(
	"log"
	// "time"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/pquerna/ffjson/ffjson"
)

type NginxLog struct {
	IPAddress   string `json:"ip_addr"`
	// TimeLocal   time   `json:"time_local"`
	RequestBody string `json:"request_body"`
	Referer     string `json:"http_referer"`
	UserAgent   string `json:"http_user_agent"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpStatus := HttpStatus{}

	defer func() { JSON(w, httpStatus) }()

	log.Println("Status command accepted")

	httpStatus.OK("Server works")
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpStatus := HttpStatus{}

	defer func() { JSON(w, httpStatus) }()

	log.Println("Push command accepted")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpStatus.ERROR("Failed to read push command body")
		return
	}

	nginxLog := &NginxLog{}
	err = ffjson.Unmarshal(body, nginxLog)
	if err != nil {
		httpStatus.ERROR("Failed to unmarshal JSON to NginxLog structure")
		return
	}

	httpStatus.OK(string(body))
}

func main() {
	server_config, err := LoadConfig("config.yml")
	if err != nil {
		log.Println("Config error: ", err)
		return
	}

	r := mux.NewRouter()

	users := r.PathPrefix("/api/v1").Subrouter()

	users.Methods("GET").
		  Path("/status").
		  HandlerFunc(statusHandler)
	users.Methods("POST").
		  Path("/push").
		  HandlerFunc(pushHandler)

	log.Println("Server running on port:", server_config.Port)
	http.ListenAndServe(":" + server_config.Port, r)
}