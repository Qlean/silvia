package main

import (
	"log"
	"net/http"
	"time"

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

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func() {
		log.Println(http.ListenAndServe(":9090", nil))
	}()
	log.Fatal(server.ListenAndServe())

	go worker.Generator()
	go worker.Transformer()
	go worker.Writer()
	go worker.Killer()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":"+worker.Config.Port, nil)
}
