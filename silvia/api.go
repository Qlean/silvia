package silvia

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/unrolled/render"
)

type Status struct {
	RabbitHealth   bool
	PostgresHealth bool
	RedshiftHealth bool

	AdjustSuccess   int
	AdjustFailed    int
	SnowplowSuccess int
	SnowplowFailed  int

	RedshiftAdjustSuccess   int
	RedshiftAdjustFailed    int
	RedshiftSnowplowSuccess int
	RedshiftSnowplowFailed  int

	PostgresAdjustSuccess   int
	PostgresAdjustFailed    int
	PostgresSnowplowSuccess int
	PostgresSnowplowFailed  int
	Uptime                  string
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

	RedshiftadjustSuccessRing := stats.RedshiftAdjustSuccessRing.Display()
	RedshiftadjustFailRing := stats.RedshiftAdjustFailRing.Display()
	RedshiftsnowplowSuccessRing := stats.RedshiftSnowplowSuccessRing.Display()
	RedshiftsnowplowFailRing := stats.RedshiftSnowplowFailRing.Display()

	PostgresadjustSuccessRing := stats.PostgresAdjustSuccessRing.Display()
	PostgresadjustFailRing := stats.PostgresAdjustFailRing.Display()
	PostgressnowplowSuccessRing := stats.PostgresSnowplowSuccessRing.Display()
	PostgressnowplowFailRing := stats.PostgresSnowplowFailRing.Display()

	status := Status{
		RabbitHealth:    stats.RabbitHealth.Get(),
		PostgresHealth:  stats.PostgresHealth.Get(),
		RedshiftHealth:  stats.RedshiftHealth.Get(),
		AdjustSuccess:   adjustSuccessRing.Total,
		AdjustFailed:    adjustFailRing.Total,
		SnowplowSuccess: snowplowSuccessRing.Total,
		SnowplowFailed:  snowplowFailRing.Total,

		RedshiftAdjustSuccess:   RedshiftadjustSuccessRing.Total,
		RedshiftAdjustFailed:    RedshiftadjustFailRing.Total,
		RedshiftSnowplowSuccess: RedshiftsnowplowSuccessRing.Total,
		RedshiftSnowplowFailed:  RedshiftsnowplowFailRing.Total,

		PostgresAdjustSuccess:   PostgresadjustSuccessRing.Total,
		PostgresAdjustFailed:    PostgresadjustFailRing.Total,
		PostgresSnowplowSuccess: PostgressnowplowSuccessRing.Total,
		PostgresSnowplowFailed:  PostgressnowplowFailRing.Total,

		Uptime: time.Since(stats.StartTime).String(),
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
	var ring []interface{}
	u, _ := url.Parse(r.URL.String())
	queryParams := u.Query()
	switch queryParams.Get("tracker") {
	case "snowplow":
		switch queryParams.Get("ring") {
		case "success":
			ring = append(ring, worker.Stats.SnowplowSuccessRing.Display(), worker.Stats.PostgresSnowplowSuccessRing.Display(), worker.Stats.RedshiftSnowplowSuccessRing.Display())
		case "failed":
			ring = append(ring, worker.Stats.SnowplowFailRing.Display(), worker.Stats.PostgresSnowplowFailRing.Display(), worker.Stats.RedshiftSnowplowFailRing.Display())
		}

	case "adjust":
		switch queryParams.Get("ring") {
		case "success":
			ring = append(ring, worker.Stats.AdjustSuccessRing.Display(), worker.Stats.PostgresAdjustSuccessRing.Display(), worker.Stats.RedshiftAdjustSuccessRing.Display())
		case "failed":
			ring = append(ring, worker.Stats.AdjustFailRing.Display(), worker.Stats.PostgresAdjustFailRing.Display(), worker.Stats.RedshiftAdjustFailRing.Display())
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
