package main

import(
	"log"
	"net/http"
)

func main() {
	var worker *Worker

	worker.Config.load()

	http.HandleFunc("/v1/status", statusHandler)
	http.HandleFunc("/v1/failed", failedHandler)
	http.HandleFunc("/v1/success", successHandler)

	go worker.generator()
	go worker.transformer()
	go worker.capacitor()

	log.Println("Server running on port:", worker.config.Port)
	http.ListenAndServe(":" + worker.config.Port, r)
}
