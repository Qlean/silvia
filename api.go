package main

import(
	"time"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/unrolled/render"
)

func statusApi(w http.ResponseWriter, r *http.Request) {
	rndr := render.New()

	counterStr   := strconv.Itoa(counter)
	failedStr    := strconv.Itoa(failRing.Total)
	uptimeStr    := time.Since(startTime).String()
	statusString := "Processed: " + counterStr + " Failed: " + failedStr + " Uptime: " + uptimeStr

	rndr.Text(w, http.StatusOK, statusString)
}

func successApi(w http.ResponseWriter, r *http.Request) {
	rndr := render.New()
	b, err := json.MarshalIndent(successRing.Ring, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	rndr.Text(w, http.StatusOK, string(b))
}

func failedApi(w http.ResponseWriter, r *http.Request) {
	rndr := render.New()
	b, err := json.MarshalIndent(failRing.Ring, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	rndr.Text(w, http.StatusOK, string(b))
}
