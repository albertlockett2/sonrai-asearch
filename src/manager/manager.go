package main

import (
	"context"
	"github.com/golang/protobuf/proto"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	queue "github.com/sonraisecurity/sonrai-asearch/src/queue"
	"log"
)

type Manager struct {
	gen.UnimplementedManagerServer
	workerQueue *queue.Queue
}

func NewManager() (*Manager, error) {
	queue, err := queue.NewQueue(queue.WORKER_QUEUE_NAME)
	if err != nil {
		// TODO log something useful here
		return nil, err
	}

	return &Manager{workerQueue: queue}, nil
}

func (m *Manager) SubmitSearch(ctx context.Context, req *gen.SubmitSearchRequest) (*gen.SubmitSearchResponse, error) {

	// TODO validate req
	// query ID not empty
	// search not empty
	// TODO validate search
	// has steps

	datas := make([][]byte, 0)

	for _, step := range req.Search.Steps {
		message := gen.InProgressRecord{
			Id: "some uuid!!",
			QueryId: req.QueryId,
			StepId: step.Id,
			PathIds: make([]*gen.RecordId, 0),
			Search: req.Search,
		}
		data, err := proto.Marshal(&message)
		if err != nil {
			return nil, err
		}
		datas = append(datas, data)
	}

	for _, data := range datas {
		log.Printf("Sending a message")
		err := m.workerQueue.Publish(data)
		if err != nil {
			return nil, err
		}
	}

	return &gen.SubmitSearchResponse{
		Status: "OK",
	}, nil
}
