// Usage, 使用方式

消费端：range Consume 方法即可

```
func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		fmt.Printf(
			"got %dB delivery: [%v] %q\n",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(true)
	}
}
```

生产端：直接调用 Publish 方法

```
