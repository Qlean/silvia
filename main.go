package main

import(
	"log"
	"time"
	"strconv"
	"net/http"
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
	statusString := "Processed: " + strconv.Itoa(counter) + " Failed: " + strconv.Itoa(failRing.Total) + " Uptime: " + time.Since(startTime).String()
	http.Error(w, statusString, http.StatusOK)
}

var geo *geoip.GeoIP
var counter int
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

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Methods("GET").
		  Path("/status").
		  HandlerFunc(statusHandler)
	api.Methods("GET").
		  Path("/failed").
		  HandlerFunc(failedHandler)
	api.Methods("GET").
		  Path("/success").
		  HandlerFunc(successHandler)

	go rabbitConsumer()

	startTime = time.Now()

	log.Println("Server running on port:", serverConfig.Port)
	http.ListenAndServe(":" + serverConfig.Port, r)
}
