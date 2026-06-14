package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewConsumer(url, queue string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Critical for controlled processing
	err = ch.Qos(
		1, // process one message at a time
		0,
		false,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // manual ack (required for reliability)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}