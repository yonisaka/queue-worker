package xasync

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	// dequeueTimeout is a default queue timeout
	dequeueTimeout = 0
)

// Queue is a struct
type Queue struct {
	client *redis.Client
	name   string
}

// NewQueue is a constructor
func NewQueue(client *redis.Client, name string) *Queue {
	return &Queue{
		client: client,
		name:   name,
	}
}

// Enqueue is a method to push item to list
func (q *Queue) Enqueue(ctx context.Context, task *Task) error {
	v, _ := json.Marshal(task)
	err := q.client.RPush(ctx, q.name, string(v)).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue: %w", err)
	}
	log.Printf("enqueued task: id=%d data=%s", task.ID, task.Data)

	return nil
}

// Dequeue is a method to pop item from list
func (q *Queue) Dequeue(ctx context.Context) (*Task, error) {
	r, err := q.client.BLPop(ctx, dequeueTimeout, q.name).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue: %w", err)
	}

	if len(r) < 2 {
		return nil, fmt.Errorf("invalid task format")
	}

	var task Task
	err = json.Unmarshal([]byte(r[1]), &task)
	if err != nil {
		return nil, fmt.Errorf("error parsing task: %w", err)
	}

	return &task, nil
}
