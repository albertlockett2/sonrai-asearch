package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/queue"
	"github.com/sonraisecurity/sonrai-asearch/src/util"
	"io/ioutil"
	"log"
	"net/http"
)

type GraphFilterMessage struct {
	Id      string
	Filters []*gen.Filter
}

type GraphEdgeMessage struct {
	Id    string
	Edges []*gen.Edge
}

type Worker struct {
	queue       *queue.Queue
	resultQueue *queue.Queue
}

func NewWorker() (*Worker, error) {
	wqueue, err := queue.NewQueue(queue.WORKER_QUEUE_NAME)
	if err != nil {
		// TODO log something useful here
		return nil, err
	}

	rqueue, err := queue.NewQueue(queue.RESULT_QUEUE_NAME)
	if err != nil {
		return nil, err
	}

	return &Worker{queue: wqueue, resultQueue: rqueue}, nil
}

func (w *Worker) Start() error {
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

func (w *Worker) Deserialize(data []byte) (*gen.InProgressRecord, error) {
	record := gen.InProgressRecord{}
	err := proto.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (w *Worker) Handle(record *gen.InProgressRecord) error {
	log.Printf("Handling message %v", record)

	step := util.FindStepById(record.StepId, record.Search)
	if step == nil {
		// TODO is this error helpful enough?
		return errors.New(fmt.Sprintf("could not find step %s in search", record.StepId))
	}

	var nextIds []*gen.RecordId
	var err error

	switch step.Type {
	case gen.SearchStep_FILTER:
		nextIds, err = w.HandleFilterStep(step, record)
	case gen.SearchStep_EDGE:
		nextIds, err = w.HandleEdgesStep(step, record)
	default:
		return errors.New(fmt.Sprintf("unknown step type %s for step %s", step.Type.String(), step.Id))
	}

	if err != nil {
		return err
	}

	for i := range nextIds {
		if len(step.NextSteps) > 0 {
			// emit next step messages
			for _, nextStep := range step.NextSteps {
				nextRecord := &gen.InProgressRecord{
					Search:  record.Search,
					StepId:  nextStep.Id,
					Id:      "some-uuid-3",
					PathIds: append(record.PathIds, nextIds[i]),
					QueryId: record.QueryId,
				}

				data, err := proto.Marshal(nextRecord)
				if err != nil {
					return err
				}
				err = w.queue.Publish(data)
				if err != nil {
					return err
				}
			}
		} else {
			// emit the result message
			resultRecord := &gen.ResultRecord{
				Id:      "some uuid",
				QueryId: record.QueryId,
				PathIds: append(record.PathIds, nextIds[i]),
				Search:  record.Search,
			}

			data, err := proto.Marshal(resultRecord)
			if err != nil {
				return err
			}

			err = w.resultQueue.Publish(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Worker) HandleFilterStep(step *gen.SearchStep, record *gen.InProgressRecord) ([]*gen.RecordId, error) {
	message := GraphFilterMessage{
		Filters: step.Filters,
	}

	if len(record.PathIds) > 0 {
		message.Id = record.PathIds[len(record.PathIds) - 1].Value
	}

	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8080/records", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return nil, err
	}

	recordIds := make([]*gen.RecordId, 0)
	for i := range ids {
		recordIds = append(recordIds, &gen.RecordId{
			StepId: step.Id,
			Value:  fmt.Sprintf("%d", ids[i]),
		})
	}
	return recordIds, nil
}

func (w *Worker) HandleEdgesStep(step *gen.SearchStep, record *gen.InProgressRecord) ([]*gen.RecordId, error) {
	// TODO
	// - validate that record has PathIds > 0
	// - error handling if request body is bad

	recordId := record.PathIds[len(record.PathIds)-1]
	message := GraphEdgeMessage{
		Id:    recordId.Value,
		Edges: step.Edges,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:8080/edges", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return nil, err
	}

	recordIds := make([]*gen.RecordId, 0)
	for i := range ids {
		recordIds = append(recordIds, &gen.RecordId{
			StepId: step.Id,
			Value:  fmt.Sprintf("%d", ids[i]),
		})
	}
	return recordIds, nil
}
