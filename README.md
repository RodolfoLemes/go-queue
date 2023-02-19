# 🛢️ Go Queue

Simple queue system supported by default array or slices and redis. Currently used to study about queue and golang. Feel free to comment, criticize or help.

## Installation

```bash
go get -u github.com/RodolfoLemes/go-queue
```

## Usage

To create a new queue or getting an existing one:
```golang
q := New[[]byte]("redis", "queueName", Options{RedisClient: rdb})
```

the T type dictated the element type inside the queue. The first argument is the driver, it can be `redis` or `array`, default is `array`. The second one is the name of the queue, the third and final one is the Options, required if you are using the redis driver.

To create a consumer, to consume the elements inside the queue, you can use:

```golang
c := NewConsumer(q)
defer c.Close()
for item := range c.Consume() {
	var data response
	json.Unmarshal(item, &data)
	log.Printf("consumer 2 got %v \n", data)
}
```

To create a producer, to append elements inside the queue, you can use:

```golang
p := NewProducer(q)
defer p.Close()

r := response{
	Data:  "b",
	IsTop: true,
	Massa: count,
}
data, _ := json.Marshal(r)

p.Produce() <- data
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)