// rabbitmq.go
package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

const (
	ExchangeName = "coffee_bed_orders"
	QueueName    = "coffee_bed_orders_queue"
	RoutingKey   = "order.created"
)

func Connect(amqpURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}

func SetupRabbitMQ(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declarar exchange
	err = channel.ExchangeDeclare(
		ExchangeName, // nombre
		"direct",     // tipo
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %v", err)
	}

	// Declarar cola
	_, err = channel.QueueDeclare(
		QueueName, // nombre
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	// Vincular cola al exchange
	err = channel.QueueBind(
		QueueName,    // nombre de la cola
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind a queue: %v", err)
	}

	return channel, nil
}

func PublishOrder(channel *amqp.Channel, order []byte) error {
	err := channel.Publish(
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        order,
		},
	)
	return err
}
