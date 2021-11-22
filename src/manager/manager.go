package manager

import (
	"context"
	"github.com/golang/protobuf/proto"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	queue "github.com/sonraisecurity/sonrai-asearch/src/queue"
	"github.com/sonraisecurity/sonrai-asearch/src/results/dao"
	"log"
)

type Manager struct {
	gen.UnimplementedManagerServer
	resultsDAO  *dao.ResultsDao
	workerQueue *queue.Queue
	// TODO not keep the history in memory
	history []*gen.SubmitSearchRequest
}

func NewManager() (*Manager, error) {
	r, err := dao.NewResultsDao()
	if err != nil {
		return nil, err
	}

	q, err := queue.NewQueue(queue.WORKER_QUEUE_NAME)
	if err != nil {
		// TODO log something useful here
		return nil, err
	}

	return &Manager{
		resultsDAO: r,
		workerQueue: q,
	}, nil
}

func (m *Manager) SubmitSearch(ctx context.Context, req *gen.SubmitSearchRequest) (*gen.SubmitSearchResponse, error) {

	// TODO validate req
	// query ID not empty
	// search not empty
	// TODO validate search
	// has steps
	m.history = append(m.history, req)

	datas := make([][]byte, 0)

	// setup tables to contain results
	err := m.resultsDAO.CreateTables(req)
	if err != nil {
		return nil, err
	}

	for _, step := range req.Search.Steps {
		message := gen.InProgressRecord{
			Id:      "some uuid!!",
			QueryId: req.QueryId,
			StepId:  step.Id,
			PathIds: make([]*gen.RecordId, 0),
			Search:  req.Search,
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

func (m *Manager) GetSearchHistory(ctx context.Context, req *gen.GetSearhHistoryRequest) (*gen.GetSearchHistoryResponse, error) {
	return &gen.GetSearchHistoryResponse{
		Requests: m.history,
	}, nil
}
