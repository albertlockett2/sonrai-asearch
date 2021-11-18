package main

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/queue"
	"io/ioutil"
	"log"
	"net/http"
)

type GraphFilterMessage struct {
	Filters []*gen.Filter
}

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

	// TODO
	// - find the actual step
	// - have different logic for processing steps
	// - take the sources and somehow add to query

	message := GraphFilterMessage{
		Filters: record.Search.Steps[0].Filters,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/records", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ids := make([]int64, 0)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return err
	}

	log.Printf("%v", ids)
	return nil
}
