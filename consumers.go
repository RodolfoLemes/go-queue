package queue

import "fmt"

type Consumer[T any] struct {
	queue *Queue[T]

	channel      chan T
	closeChannel chan struct{}
	errChannel   chan error
}

func NewConsumer[T any](q *Queue[T]) *Consumer[T] {
	return &Consumer[T]{
		queue:        q,
		channel:      make(chan T),
		closeChannel: make(chan struct{}),
		errChannel:   make(chan error),
	}
}

func (c *Consumer[T]) Consume() <-chan T {
	go func() {
		for {
			select {
			default:
				item, err := c.queue.dequeue()
				if err != nil {
					c.errChannel <- fmt.Errorf("queue error: %w", err)
				}

				c.channel <- item
			case <-c.closeChannel:
				return
			}
		}
	}()

	return c.channel
}

func (c *Consumer[T]) Close() error {
	<-c.closeChannel

	close(c.channel)
	close(c.closeChannel)
	return nil
}

func (c *Consumer[T]) Err() <-chan error {
	return c.errChannel
}
