package main

import(
	"log"
	"time"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/abh/geoip"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Status command accepted")
	statusString := "Processed: " + strconv.Itoa(counter) + " Failed: " + strconv.Itoa(failCounter) + " Uptime: " + time.Since(startTime).String()
	http.Error(w, statusString, http.StatusOK)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	localCounter := counter

	log.Println(counter, "Push command accepted")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		eventFail(w, localCounter, "Failed to read push command body", err)
		return
	}

	log.Println(counter, "Raw event:", string(body))

	// Remove Line Feed
	body = body[0:len(body)-1]

	normString, err := strconv.Unquote(`"` + string(body) + `"`)
	if err != nil {
		eventFail(w, localCounter, "Failed to normalize nginx log string", err)
		return
	}

	nginxLog := &NginxLog{}
	err = json.Unmarshal([]byte(normString), nginxLog)
	if err != nil {
		eventFail(w, localCounter, "Failed to unmarshal JSON to NginxLog structure", err)
		return
	}

	err = dbmap.Insert(nginxLogToEvent(nginxLog))
	if err != nil {
		eventFail(w, localCounter, "Failed to write data to database", err)
		return
	}

	log.Println(counter, "Succesfuly witten")
}

var geo *geoip.GeoIP
var counter int
var failCounter int
var startTime time.Time

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

	startTime = time.Now()

	log.Println("Server running on port:", serverConfig.Port)
	http.ListenAndServe(":" + serverConfig.Port, r)
}
