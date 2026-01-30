package redis

import (
	"context"
	"log"
	"os"

	"guesswhoclientapi/sse"

	"github.com/go-redis/redis/v8"	
)

var ctx = context.Background()

// SubscribeToUpdates subscribes to the "game_updates" channel in Redis
// and broadcasts an "update" message to the SSE broker.
func SubscribeToUpdates(broker *sse.Broker) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		// If REDIS_ADDR is not set, we don't start the subscriber.
		return
	}

	log.Println("Redis address found, initializing Redis client...")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Check if the connection is successful
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Could not connect to Redis: %v. Real-time updates will not be available.", err)
		return
	}

	log.Println("Successfully connected to Redis. Subscribing to 'game_updates' channel.")

	pubsub := rdb.Subscribe(ctx, "game_updates")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Received message from Redis on channel '%s'. Triggering SSE update.", msg.Channel)
		// Broadcast a generic "update" message, same as the HTTP endpoint.
		broker.Broadcast("update")
	}
}