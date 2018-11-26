package main

import (
	"log"
	"net/http"

	"github.com/Qlean/silvia/silvia"

	_ "net/http/pprof"
)

func main() {
	worker := &silvia.Worker{}
	err := worker.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/v1/status", worker.ApiHandler(silvia.StatusApi))
	http.Handle("/v1/ring", worker.ApiHandler(silvia.RingApi))
	http.Handle("/debug/pprof/", pprof.Index)
	http.Handle("/debug/pprof/cmdline", pprof.Cmdline)
	http.Handle("/debug/pprof/profile", pprof.Profile)
	http.Handle("/debug/pprof/symbol", pprof.Symbol)
	http.Handle("/debug/pprof/trace", pprof.Trace)

	go worker.Generator()
	go worker.Transformer()
	go worker.Writer()
	go worker.Killer()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":"+worker.Config.Port, nil)
}
