package main

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"

	. "github.com/RodolfoLemes/go-queue"
)

func main() {
	// arr := make(Array[int], 0)
	type response struct {
		Data  string
		IsTop bool
		Massa int
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	q := New[[]byte]("redis", "queueName", Options{RedisClient: rdb})

	q.Flush()

	go func() {
		c := NewConsumer(q)
		defer c.Close()

		for item := range c.Consume() {
			var data response
			json.Unmarshal(item, &data)
			log.Printf("consumer 1 got %v \n", data)
		}
	}()
	go func() {
		c := NewConsumer(q)
		defer c.Close()

		for item := range c.Consume() {
			var data response
			json.Unmarshal(item, &data)
			log.Printf("consumer 2 got %v \n", data)
		}
	}()
	go func() {
		c := NewConsumer(q)
		defer c.Close()

		for item := range c.Consume() {
			var data response
			json.Unmarshal(item, &data)
			log.Printf("consumer 3 got %v \n", data)
		}
	}()

	count := 0
	p := NewProducer(q)
	defer p.Close()
	for {
		r := response{
			Data:  "b",
			IsTop: true,
			Massa: count,
		}

		data, _ := json.Marshal(r)

		p.Produce() <- data

		count++
	}
}
