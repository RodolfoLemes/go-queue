package queue

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Options struct {
	RedisClient *redis.Client
}

type Queue[T any] struct {
	driver list[T]

	queueName string
	capacity  int64

	m *sync.Mutex
}

func New[T any](driverName string, queueName string, opts Options) *Queue[T] {
	var driver list[T]
	switch driverName {
	default:
		fallthrough
	case "array":
		driver = &array[T]{}
	case "redis":
		driver = newRedisList[T](queueName, opts.RedisClient)
	}

	return &Queue[T]{
		driver:   driver,
		capacity: 1000,
		m:        &sync.Mutex{},
	}
}

func (q *Queue[T]) enqueue(data T) error {
	q.m.Lock()
	defer q.m.Unlock()
	if q.capacity == q.driver.Size() {
		return fmt.Errorf("queue is full")
	}

	return q.driver.Append(data)
}

func (q *Queue[T]) dequeue() (item T, err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.driver.Size() == 0 {
		return item, fmt.Errorf("queue is empty")
	}

	return q.driver.Detach()
}

func (q *Queue[T]) Flush() error {
	q.m.Lock()
	defer q.m.Unlock()
	return q.driver.Flush()
}
