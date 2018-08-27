package silvia

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/unrolled/render"
)

type Status struct {
	RabbitHealth    bool
	PostgresHealth  bool
	AdjustSuccess   int
	AdjustFailed    int
	SnowplowSuccess int
	SnowplowFailed  int
	Uptime          string
}

func (worker *Worker) ApiHandler(fn func(http.ResponseWriter, *http.Request, *Worker)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, worker)
	}
}

func StatusApi(w http.ResponseWriter, r *http.Request, worker *Worker) {
	rndr := render.New()

	stats := worker.Stats

	adjustSuccessRing := stats.AdjustSuccessRing.Display()
	adjustFailRing := stats.AdjustFailRing.Display()
	snowplowSuccessRing := stats.SnowplowSuccessRing.Display()
	snowplowFailRing := stats.SnowplowFailRing.Display()

	status := Status{
		RabbitHealth:    stats.RabbitHealth.Get(),
		PostgresHealth:  stats.PostgresHealth.Get(),
		AdjustSuccess:   adjustSuccessRing.Total,
		AdjustFailed:    adjustFailRing.Total,
		SnowplowSuccess: snowplowSuccessRing.Total,
		SnowplowFailed:  snowplowFailRing.Total,
		Uptime:          time.Since(stats.StartTime).String(),
	}

	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	httpStatus := http.StatusOK

	if !status.RabbitHealth || !status.PostgresHealth {
		httpStatus = http.StatusTooManyRequests
	}

	rndr.Text(w, httpStatus, string(b))
}

func RingApi(w http.ResponseWriter, r *http.Request, worker *Worker) {
	var ring interface{}
	u, _ := url.Parse(r.URL.String())
	queryParams := u.Query()
	switch queryParams.Get("tracker") {
	case "snowplow":
		switch queryParams.Get("ring") {
		case "success":
			ring = worker.Stats.SnowplowSuccessRing.Display()
		case "failed":
			ring = worker.Stats.SnowplowFailRing.Display()
		}
	case "adjust":
		switch queryParams.Get("ring") {
		case "success":
			ring = worker.Stats.AdjustSuccessRing.Display()
		case "failed":
			ring = worker.Stats.AdjustFailRing.Display()
		}
	}

	rndr := render.New()
	b, err := json.MarshalIndent(ring, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	rndr.Text(w, http.StatusOK, string(b))
}
