package consumers

import (
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume()
}

type ConsumerImpl struct {
	q  *amqp091.Queue
	ch *amqp091.Channel
	Wg sync.WaitGroup
}

func (c *ConsumerImpl) Consume() error {
	msgs, err := c.ch.Consume(
		c.q.Name,
		"",
		true,  // Auto-acknowledge
		false, // Exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

}
func (c *ConsumerImpl) worker(id int, msgs <-chan amqp091.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range msgs {
		fmt.Printf("👷 Worker %d processing: %s\n", id, msg.Body)

		// Simulate work (sleep for 2 seconds)
		time.Sleep(2 * time.Second)

		// Acknowledge message
		msg.Ack(false)
		fmt.Printf("✅ Worker %d finished processing: %s\n", id, msg.Body)
	}
}
