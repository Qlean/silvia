package main

import(
	"log"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/abh/geoip"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Status command accepted")
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Push command accepted")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read push command body: " + err.Error(), http.StatusBadRequest)
		log.Println("Failed to read push command body: " + err.Error())
		return
	}

	normString, err := strconv.Unquote(`"` + string(body) + `"`)
	if err != nil {
		http.Error(w, "Failed to normalize nginx log string: " + err.Error(), http.StatusBadRequest)
		log.Println("Failed to normalize nginx log string: " + err.Error())
		return
	}

	nginxLog := &NginxLog{}
	err = json.Unmarshal([]byte(normString), nginxLog)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON to NginxLog structure: " + err.Error(), http.StatusBadRequest)
		log.Println("Failed to unmarshal JSON to NginxLog structure: " + err.Error())
		return
	}

	log.Println(nginxLog)

	err = dbmap.Insert(nginxLogToEvent(nginxLog))
    if err != nil {
		http.Error(w, "Failed to write data to database: " + err.Error(), http.StatusBadRequest)
		log.Println("Failed to write data to database: " + err.Error())
		return
	}
}

var geo *geoip.GeoIP

func main() {
	var err error
	serverConfig, err = LoadConfig("config.yaml")
	if err != nil {
		log.Println("Config error: ", err)
		return
	}

	geo, err = geoip.Open("GeoLiteCity.dat")
	if err != nil {
		log.Println("Could not open GeoIPv6 database: %s\n", err)
	}

	dbmap = initDb()
    defer dbmap.Db.Close()

	r := mux.NewRouter()

	users := r.PathPrefix("/api/v1").Subrouter()

	users.Methods("GET").
		  Path("/status").
		  HandlerFunc(statusHandler)
	users.Methods("POST").
		  Path("/push").
		  HandlerFunc(pushHandler)

	log.Println("Server running on port:", serverConfig.Port)
	http.ListenAndServe(":" + serverConfig.Port, r)
}
