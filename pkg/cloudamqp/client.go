package cloudamqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    string
	hostname string
	done     chan struct{} // 用于控制消费者goroutine的关闭
}

// NewClient creates a new CloudAMQP client
func NewClient(uri string, queue string, hostname string) (*Client, error) {
	// Create the connection
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	client := &Client{
		conn:     conn,
		channel:  ch,
		queue:    queue,
		hostname: hostname,
		done:     make(chan struct{}),
	}

	// 启动消费者
	go client.startConsumer()

	return client, nil
}

// startConsumer 启动一个消费者来处理队列中的消息
func (c *Client) startConsumer() {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Printf("Failed to start consumer: %v", err)
		return
	}

	for {
		select {
		case <-c.done:
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			// 这里我们只是记录收到的消息
			log.Printf("Received keep-alive message: %s", string(msg.Body))
		}
	}
}

// Ping sends a keep-alive message to the queue
func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create keep-alive message
	message := struct {
		Timestamp time.Time `json:"ping_timestamp"`
		Source    string    `json:"ping_source"`
		Details   struct {
			Hostname string `json:"hostname"`
			Version  string `json:"version"`
		} `json:"ping_details"`
	}{
		Timestamp: time.Now(),
		Source:    "cloudamqp-keeper",
		Details: struct {
			Hostname string `json:"hostname"`
			Version  string `json:"version"`
		}{
			Hostname: c.hostname,
			Version:  "1.0",
		},
	}

	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish the message
	err = c.channel.PublishWithContext(ctx,
		"",       // exchange
		c.queue,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close closes the RabbitMQ connection and channel
func (c *Client) Close() error {
	// 首先关闭消费者
	close(c.done)

	var err error
	if c.channel != nil {
		if err = c.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
} 