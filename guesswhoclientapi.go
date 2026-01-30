package main

import (
	"log"
	"net/http"

	"guesswhoclientapi/api"
	"guesswhoclientapi/redis"
)

func main() {
	router, broker := api.NewRouterAndBroker()

	// Start the Redis subscriber in a background goroutine
	go redis.SubscribeToUpdates(broker)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}