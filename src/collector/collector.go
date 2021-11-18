package main

import (
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/queue"
	"google.golang.org/protobuf/proto"
	"log"
)

type Collector struct {
	queue *queue.Queue
}

func NewCollector() (*Collector, error) {
	q, err := queue.NewQueue(queue.RESULT_QUEUE_NAME)
	if err != nil {
		return nil, err
	}
	return &Collector{
		queue: q,
	}, nil
}

func (c *Collector) Start() error {
	msgs, err := c.queue.Consume()
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message")
			record, err := c.Deserialize(d.Body)
			if err != nil {
				log.Printf("Error deserializing %v", err)
				continue
			}

			log.Printf("recieved_record: %v", record.PathIds)

			//err = c.Handle(record)
			//if err != nil {
			//	log.Printf("Error handling %v", err)
			//	continue
			//}
		}
	}()
	<-forever

	return nil
}

func (c *Collector) Deserialize(data []byte) (*gen.ResultRecord, error) {
	record := gen.ResultRecord{}
	err := proto.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
