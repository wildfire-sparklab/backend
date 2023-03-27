package rabbit

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type Conn struct {
	Channel *amqp091.Channel
}

// GetConn -
func GetConn(rabbitURL string) (Conn, error) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return Conn{}, err
	}

	ch, err := conn.Channel()
	_, err = ch.QueueDeclare(
		"broker", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	return Conn{
		Channel: ch,
	}, err
}

func (conn Conn) Publish(routingKey string, data []byte) error {
	return conn.Channel.PublishWithContext(context.TODO(),
		// exchange - yours may be different
		"",
		routingKey,
		// mandatory - we don't care if there I no queue
		false,
		// immediate - we don't care if there is no consumer on the queue
		false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: amqp091.Persistent,
		})
}
