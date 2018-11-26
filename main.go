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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go worker.Generator()
	go worker.Transformer()
	go worker.Writer()
	go worker.Killer()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":"+worker.Config.Port, nil)
}
