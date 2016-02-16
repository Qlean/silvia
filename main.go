package main

import(
	"log"
	"net/http"
	"./silvia"
)

func main() {
	var worker *silvia.Worker

	err := worker.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/v1/status", statusHandler)
	http.HandleFunc("/v1/success", successHandler)
	http.HandleFunc("/v1/failed", failedHandler)

	go worker.Generator()
	go worker.Transformer()
	go worker.Capacitor()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":" + worker.Config.Port, nil)
}
