package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	Limit     = 5
	TimeLimit = 10
)

var (
	ctx    = context.Background()
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

func main() {
	defer client.Close()

	for i := 0; i < 10; i++ {
		isAllowed, err := checkRateLimit("user1")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if isAllowed {
			fmt.Println("Allowed %d is allowed\n", i+1)
		} else {
			fmt.Println("Not Allowed\n", i+1)
		}
		time.Sleep(1 * time.Second)
	}
}

func checkRateLimit(userID string) (bool, error) {
	key := fmt.Sprintf("hit:%s", userID)
	// Use a transaction to increment the counter and set the expiry
	pipe := client.TxPipeline()
	// increment the counter
	pipe.Incr(key)

	// Set the expiry if the key is new
	pipe.Expire(key, TimeLimit*time.Second)
	_, err := pipe.Exec()

	if err != nil {
		return false, err
	}
	count, err := client.Get(key).Int()
	if err != nil {
		return false, err
	}

	return count <= Limit, nil
}
