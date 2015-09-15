package main

import(
	"log"
	"time"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/pquerna/ffjson/ffjson"
)

type NginxLog struct {
	IPAddress   string    `json:"ip_addr"`
	TimeLocal   time.Time `json:"time_local"`
	RequestBody string    `json:"request_body"`
	Referer     string    `json:"http_referer"`
	UserAgent   string    `json:"http_user_agent"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Status command accepted")
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Push command accepted")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read push command body", http.StatusBadRequest)
		return
	}

	normString, err := strconv.Unquote(`"` + string(body) + `"`)
	if err != nil {
		http.Error(w, "Failed to normalize nginx log string: " + err.Error(), http.StatusBadRequest)
		return
	}

	nginxLog := &NginxLog{}
	err = ffjson.Unmarshal([]byte(normString), nginxLog)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON to NginxLog structure: " + err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	server_config, err := LoadConfig("config.yaml")
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