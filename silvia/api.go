package silvia

import(
	"time"
	"net/url"
	"net/http"
	"encoding/json"

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
	status := Status{
		RabbitHealth:    stats.RabbitHealth,
		PostgresHealth:  stats.PostgresHealth,
		AdjustSuccess:   stats.AdjustSuccessRing.Total,
		AdjustFailed:    stats.AdjustFailRing.Total,
		SnowplowSuccess: stats.SnowplowSuccessRing.Total,
		SnowplowFailed:  stats.SnowplowFailRing.Total,
		Uptime:          time.Since(stats.StartTime).String(),
	}

	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		rndr.Text(w, http.StatusBadRequest, "Cant draw pretty JSON")
		return
	}

	httpStatus := http.StatusOK

	if !status.RabbitHealth || !status.PostgresHealth { httpStatus = http.StatusTooManyRequests }

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
			ring = worker.Stats.SnowplowSuccessRing
		case "failed":
			ring = worker.Stats.SnowplowFailRing
		}
	case "adjust":
		switch queryParams.Get("ring") {
			case "success":
				ring = worker.Stats.AdjustSuccessRing
			case "failed":
				ring = worker.Stats.AdjustFailRing
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
