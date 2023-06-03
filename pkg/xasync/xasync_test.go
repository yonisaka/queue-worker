package xasync_test

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/queue-worker/pkg/xasync"
	"github.com/yonisaka/queue-worker/tests"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestQueue_Async tests the New function of the redis queue
func TestQueue_Async(t *testing.T) {
	var (
		testbox   = tests.Init()
		cfg       = testbox.Cfg
		ctx       = testbox.Ctx
		queueName = "task:delivery"
	)

	// Create a Redis client
	dns := fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     dns,
		Password: cfg.RedisConfig.Password, // Empty if no password is set
		DB:       cfg.RedisConfig.DB,       // Default database
	})

	// Ping the Redis server to ensure connectivity
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Create the task queue
	taskQueue := xasync.NewQueue(client, queueName)

	_ = taskQueue.Enqueue(ctx, &xasync.Task{
		ID:   1,
		Data: "Hi, task 1",
	})
	_ = taskQueue.Enqueue(ctx, &xasync.Task{
		ID:   2,
		Data: "Hi, task 2",
	})
	_ = taskQueue.Enqueue(ctx, &xasync.Task{
		ID:   3,
		Data: "Hi, task 3",
	})
	_ = taskQueue.Enqueue(ctx, &xasync.Task{
		ID:   4,
		Data: "Hi, task 4",
	})

	assert.True(t, true)
}

func TestWorker_Async(t *testing.T) {
	var (
		testbox   = tests.Init()
		cfg       = testbox.Cfg
		ctx       = testbox.Ctx
		queueName = "task:delivery"
	)

	// Create a Redis client
	dns := fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     dns,
		Password: cfg.RedisConfig.Password, // Empty if no password is set
		DB:       cfg.RedisConfig.DB,       // Default database
	})

	// Ping the Redis server to ensure connectivity
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	taskQueue := xasync.NewQueue(client, queueName)

	// Create a wait group to ensure all worker goroutines finish
	var wg sync.WaitGroup

	// Create a channel to receive termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the workers
	worker := xasync.NewWorker(2, taskQueue, &wg)
	worker.StartWorkers()

	// Wait for termination signal or completion of all tasks
	select {
	case <-sigChan:
		// Termination signal received, initiate graceful shutdown
		fmt.Println("Termination signal received. Initiating graceful shutdown...")
		closeSigChan := make(chan struct{})
		go func() {
			close(closeSigChan)
		}()

		select {
		case <-closeSigChan:
			fmt.Println("All tasks processed. Exiting...")
		case <-time.After(5 * time.Second):
			fmt.Println("Graceful shutdown timed out. Exiting...")
		}
	}
}
