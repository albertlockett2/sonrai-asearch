package main

import (
	"github.com/golang/protobuf/proto"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/queue"
	"log"
)

type Worker struct {
	queue *queue.Queue
}

func NewWorker() (*Worker, error) {
	q, err := queue.NewQueue()
	if err != nil {
		// TODO log something useful here
		return nil, err
	}

	return &Worker{queue: q}, nil
}

func (w*Worker) Start() error {
	msgs, err := w.queue.Consume()
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message")
			record, err := w.Deserialize(d.Body)
			if err != nil {
				log.Printf("Error deserializing %v", err)
				continue
			}

			err = w.Handle(record)
			if err != nil {
				log.Printf("Error handling %v", err)
				continue
			}
		}
	}()
	<-forever

	return nil
}

func (w*Worker) Deserialize(data []byte) (*gen.InProgressRecord, error) {
	record := gen.InProgressRecord{}
	err := proto.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (w*Worker) Handle(record *gen.InProgressRecord) error {
	log.Printf("Handling message %v", record)
	return nil
}
