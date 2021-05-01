package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	Client 		*redis.Client
	Context 	context.Context
}

// NewClient connects to redis server
// and returns a pointer to a RedisClient and nil if no error
func NewClient(addr string) (*RedisClient, error) {
	ctx := context.TODO()

	// Connect to redis server
	client := redis.NewClient(&redis.Options{
		Addr:     addr, // e.g. "localhost:6379"
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// If Ping throws error
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: client,
		Context: ctx,
	}, nil
}

// Set calls the redis set command
func (r *RedisClient) Set(key, value string, expTime time.Duration, numReplicas int) error {
	err := r.Client.Set(r.Context, key, value, expTime).Err()
	if err != nil {
		return err
	}

	// Use redis WAIT command to block until numReplicas replicas got the previous write
	n, err := r.Client.Do(r.Context, "wait", numReplicas, 0).Result()
	fmt.Println(n)

	return nil
}

// Get calls the redis get command
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.Client.Get(r.Context, key).Result()

	switch {
	case err == redis.Nil:
		return "", errors.New("key does not exist")
	case err != nil:
		return "", err
	case val == "":
		return "", errors.New("value is empty")
	default:
		return val, nil
	}
}

func (r *RedisClient) PrintInfo() {
	fmt.Println(r.Client.Do(r.Context, "info", "replication"))
}
