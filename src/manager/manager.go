package main

import (
	"context"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	queue "github.com/sonraisecurity/sonrai-asearch/src/queue"
)

type Manager struct {
	gen.UnimplementedManagerServer
	queue *queue.Queue
}

func NewManager() (*Manager, error) {
	queue, err := queue.NewQueue()
	if err != nil {
		// TODO log something useful here
		return nil, err
	}

	return &Manager{queue: queue}, nil
}

func (m *Manager) SubmitSearch(ctx context.Context, req *gen.SubmitSearchRequest) (*gen.SubmitSearchResponse, error) {
	err := m.queue.Publish()
	if err != nil {
		return nil, err
	}

	return &gen.SubmitSearchResponse{
		Status: "OK",
	}, nil
}
