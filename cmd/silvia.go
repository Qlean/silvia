package main

import (
	"log"
	"net/http"

	"github.com/Qlean/silvia/silvia"
)

func main() {
	worker := &silvia.Worker{}
	err := worker.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/v1/status", worker.ApiHandler(silvia.StatusApi))
	http.Handle("/v1/ring", worker.ApiHandler(silvia.RingApi))

	go worker.Generator()
	go worker.Transformer()
	go worker.Writer("postgres")
	go worker.Writer("redshift")
	go worker.Killer()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":"+worker.Config.Port, nil)
}
