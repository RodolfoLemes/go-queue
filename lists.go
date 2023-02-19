package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type list[T any] interface {
	Append(T) error
	Detach() (T, error)
	Size() int64
	Flush() error
}

type array[T any] []T

var _ list[interface{}] = (*array[interface{}])(nil)

func (a *array[T]) Append(item T) error {
	*a = append(*a, item)
	return nil
}

func (a *array[T]) Detach() (item T, err error) {
	aux := make([]T, len(*a)-1)

	for i := range *a {
		if i == 0 {
			item = (*a)[i]
			continue
		}

		aux[i-1] = (*a)[i]
	}

	*a = aux

	return
}

func (a *array[T]) Size() int64 {
	return int64(len(*a))
}

func (a *array[T]) Flush() error {
	*a = array[T]{}
	return nil
}

type redisList[T any] struct {
	ctx    context.Context
	key    string
	client *redis.Client
}

var _ list[interface{}] = (*redisList[interface{}])(nil)

func newRedisList[T any](key string, client *redis.Client) *redisList[T] {
	return &redisList[T]{
		ctx:    context.Background(),
		key:    key,
		client: client,
	}
}

func (r *redisList[T]) Append(item T) error {
	err := r.client.LPush(
		r.ctx,
		r.key,
		item,
	).Err()

	return err
}

func (r *redisList[T]) Detach() (item T, err error) {
	err = r.client.RPop(
		r.ctx,
		r.key,
	).Scan(&item)

	return item, err
}

func (r *redisList[T]) Size() int64 {
	size := r.client.LLen(
		r.ctx,
		r.key,
	).Val()

	return size
}

func (r *redisList[T]) Flush() error {
	err := r.client.Del(
		r.ctx,
		r.key,
	).Err()

	return err
}
