package queue

import (
	"fmt"
)

type Producer[T any] struct {
	queue *Queue[T]

	channel      chan T
	closeChannel chan struct{}
	errChannel   chan error
}

func NewProducer[T any](q *Queue[T]) *Producer[T] {
	channel := make(chan T)
	closeChannel := make(chan struct{})
	errChannel := make(chan error)

	go func() {
		for {
			select {
			case value := <-channel:
				err := q.enqueue(value)
				if err != nil {
					errChannel <- fmt.Errorf("queue error: %w", err)
				}
			case <-closeChannel:
				return
			}
		}
	}()

	return &Producer[T]{
		queue:        q,
		channel:      channel,
		closeChannel: closeChannel,
	}
}

func (c *Producer[T]) Produce() chan<- T {
	return c.channel
}

func (c *Producer[T]) Close() error {
	<-c.closeChannel

	close(c.channel)
	close(c.closeChannel)
	return nil
}

func (c *Producer[T]) Err() <-chan error {
	return c.errChannel
}
