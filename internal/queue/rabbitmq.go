package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const CIQueueName = "CIJobs"

type RMQQueue struct {
	rmqConnection *amqp.Connection
	rmqChannel    *amqp.Channel

	CIJobsQueue amqp.Queue
}

func NewRMQQueue(addr string) *RMQQueue {
	conn, err := amqp.Dial(addr)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open RabbitMQ channel")

	ciQueue, err := declareQueue(ch, CIQueueName)
	failOnError(err, "Failed to create CI queue")

	return &RMQQueue{
		rmqConnection: conn,
		rmqChannel:    ch,
		CIJobsQueue:   ciQueue,
	}
}

func (q *RMQQueue) MakeCIMsgChan() (<-chan []byte, error) {
	msgs, err := q.rmqChannel.Consume(
		CIQueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgsChan := make(chan []byte)
	go func() {
		for m := range msgs {
			msgsChan <- m.Body
		}
	}()

	return msgsChan, nil
}

func (q *RMQQueue) Close() {
	q.rmqConnection.Close()
}

func declareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
