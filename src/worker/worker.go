package main

import (
	"bytes"
	"github.com/sonraisecurity/sonrai-asearch/src/queue"
	"log"
	"time"
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
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	}()
	<-forever

	return nil
}

