package xasync_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/yonisaka/queue-worker/pkg/xasync"
	"github.com/yonisaka/queue-worker/tests"
	"log"
	"testing"
)

// TestQueue_Hibiken tests the New function of queue hibiken/asynq package
func TestQueue_Hibiken(t *testing.T) {
	var (
		testbox   = tests.Init()
		cfg       = testbox.Cfg
		queueName = "task:delivery"
	)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port)})
	defer client.Close()

	task1, err := newDeliveryTask(queueName, 1, "Hi, task 1")
	if err != nil {
		t.Error(err)
	}

	info, err := client.Enqueue(task1)
	if err != nil {
		t.Error(err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	task2, err := newDeliveryTask(queueName, 2, "Hi, task 2")
	if err != nil {
		t.Error(err)
	}

	info, err = client.Enqueue(task2)
	if err != nil {
		t.Error(err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	assert.True(t, true)
}

// TestWorker_Hibiken tests the New function of worker hibiken/asynq package
func TestWorker_Hibiken(t *testing.T) {
	var (
		testbox   = tests.Init()
		cfg       = testbox.Cfg
		queueName = "task:delivery"
	)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port)},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 1,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queueName, handleDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
	assert.True(t, true)
}

func newDeliveryTask(queueName string, ID int, data string) (*asynq.Task, error) {
	task := xasync.Task{
		ID:   ID,
		Data: data,
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(queueName, payload), nil
}

func handleDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p xasync.Task
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Println(fmt.Sprintf("handle delivery task id=%d data=%s", p.ID, p.Data))

	return nil
}
