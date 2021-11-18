package queue

import "github.com/streadway/amqp"

const WORKER_QUEUE_NAME string = "q_asearch_worker"
const RESULT_QUEUE_NAME string = "q_asearch_results"

type Queue struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

func NewQueue(name string) (*Queue, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close()

	queue, err := ch.QueueDeclare(
		name,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		conn:  conn,
		ch:    ch,
		queue: &queue,
	}, nil
}

func (q *Queue) Publish(data []byte) error {
	err := q.ch.Publish(
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)

	return err
}

func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := q.ch.Consume(
		q.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
