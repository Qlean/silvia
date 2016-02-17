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

	http.HandleFunc("/v1/status", statusApi)
	http.HandleFunc("/v1/success", successApi)
	http.HandleFunc("/v1/failed", failedApi)

	go worker.Generator()
	go worker.Transformer()
	go worker.Writer()

	log.Println("Server running on port:", worker.Config.Port)
	http.ListenAndServe(":" + worker.Config.Port, nil)
}
