package broker

import (
	"context"
	"encoding/json"
	"log"
	"telemetry-service/internal/processor"

	"github.com/redis/go-redis/v9"
)

// RedisBroker handles publishing processed telemetry to Redis channels.
type RedisBroker struct {
	client *redis.Client
}

func NewRedisBroker(addr string, password string, db int) *RedisBroker {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	// Ping to test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("[WARNING] Redis connection failed: %v", err)
	} else {
		log.Printf("Redis connected at %s", addr)
	}

	return &RedisBroker{client: rdb}
}

// Publish broadcasts the LocomotiveState to the specified Redis channel.
func (rb *RedisBroker) Publish(ctx context.Context, channel string, state *processor.LocomotiveState) error {
	payload, err := json.Marshal(state)
	if err != nil {
		return err
	}
	// The Laravel app can subscribe to this channel using `Illuminate\Support\Facades\Redis`
	return rb.client.Publish(ctx, channel, payload).Err()
}
