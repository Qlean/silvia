package main

import(
	"log"
	"net/http"
)

func main() {
	var worker *Worker

	err := worker.load()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/v1/status", statusHandler)
	http.HandleFunc("/v1/success", successHandler)
	http.HandleFunc("/v1/failed", failedHandler)

	go worker.generator()
	go worker.transformer()
	go worker.capacitor()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":" + worker.Config.Port, nil)
}
