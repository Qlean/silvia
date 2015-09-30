package main

import(
	"log"
	"time"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/abh/geoip"
)

func successHandler(w http.ResponseWriter, r *http.Request) {
	rndr := render.New()
	b, err := json.MarshalIndent(successRing.Ring, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	rndr.Text(w, http.StatusOK, string(b))
}

func failedHandler(w http.ResponseWriter, r *http.Request) {
	rndr := render.New()
	b, err := json.MarshalIndent(failRing.Ring, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	rndr.Text(w, http.StatusOK, string(b))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Status command accepted")
	statusString := "Processed: " + strconv.Itoa(counter) + " Failed: " + strconv.Itoa(failCounter) + " Uptime: " + time.Since(startTime).String()
	http.Error(w, statusString, http.StatusOK)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	localCounter := counter

	var err error
	var body []byte
	var nginxLog *NginxLog = &NginxLog{}
	var event *Event

	log.Println(counter, "Push command accepted")

	defer r.Body.Close()
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		eventFail(w, localCounter, "Failed to read push command body", err, "", nginxLog, event)
		return
	}

	// Remove Line Feed
	body = body[0:len(body)-1]
	rawBody := string(body)

	normString, err := strconv.Unquote(`"` + rawBody + `"`)
	if err != nil {
		eventFail(w, localCounter, "Failed to normalize nginx log string", err, rawBody, nginxLog, event)
		return
	}

	err = json.Unmarshal([]byte(normString), nginxLog)
	if err != nil {
		eventFail(w, localCounter, "Failed to unmarshal JSON to NginxLog structure", err, normString, nginxLog, event)
		return
	}

	event = nginxLogToEvent(nginxLog)
	err = dbmap.Insert(event)
	if err != nil {
		eventFail(w, localCounter, "Failed to write data to database", err, normString, nginxLog, event)
		return
	}

	successRing.Add(normString, nginxLog, event, "success")
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
	users.Methods("GET").
		  Path("/failed").
		  HandlerFunc(failedHandler)
	users.Methods("GET").
		  Path("/success").
		  HandlerFunc(successHandler)
	users.Methods("POST").
		  Path("/push").
		  HandlerFunc(pushHandler)

	startTime = time.Now()

	log.Println("Server running on port:", serverConfig.Port)
	http.ListenAndServe(":" + serverConfig.Port, r)
}
