package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/util"
	"log"
	"strings"
)

type resultRecord struct {
	StepId   string
	RecordId string
	OtherId  string
}

func (r *resultRecord) documentId() string {
	return fmt.Sprintf("%s_%s_%s", r.StepId, r.RecordId, r.OtherId)
}

type ResultsDao struct {
	esclient *elasticsearch.Client
}

func NewResultsDao() (*ResultsDao, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:19200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ResultsDao{
		esclient: es,
	}, nil
}

func (r *ResultsDao) CreateTables(req *gen.SubmitSearchRequest) error {
	createIndexReq := esapi.IndicesCreateRequest{
		Index: indexName(req.QueryId),
		Body: strings.NewReader(`{
			"mappings": {
				"properties": {
					"StepId": { "type": "keyword" },
					"RecordId": { "type": "keyword" },
					"OtherId": { "type": "keyword" }
				}
      }
		}`),
	}

	res, err := createIndexReq.Do(context.Background(), r.esclient)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting response: %s", err))
	}
	defer res.Body.Close()

	// TODO check if body has error in it
	return nil
}

func (r *ResultsDao) WriteResults(results *gen.ResultRecord) error {
	stepsById := make(map[string]*gen.SearchStep) // cache of steps

	for i, result := range results.PathIds {
		// find step, add to cache
		if stepsById[result.StepId] == nil {
			step := util.FindStepById(result.StepId, results.Search)
			if step == nil {
				return errors.New(fmt.Sprintf("could not find step for id %s", result.StepId))
			}
			stepsById[result.StepId] = step
		}

		// construct record
		record := &resultRecord{StepId: result.StepId}
		if stepsById[result.StepId].Type == gen.SearchStep_EDGE {
			record.RecordId = results.PathIds[i-1].Value
			record.OtherId = result.Value
		} else {
			record.RecordId = result.Value
		}

		// write record
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}

		req := esapi.IndexRequest{
			Index:      indexName(results.QueryId),
			DocumentID: record.documentId(),
			Body:       strings.NewReader(string(data)),
			Refresh:    "false",
			OpType:     "create",
		}

		res, err := req.Do(context.Background(), r.esclient)
		if err != nil {
			return errors.New(fmt.Sprintf("error getting response: %s", err))
		}
		defer res.Body.Close()

	}

	return nil
}

func (r *ResultsDao) GetResults(req *gen.ResultsRequest) ([]string, error) {
	// TODO
	// - don't build the search body like this
	// - rewrite this method so it's not so weird

	searchReq := esapi.SearchRequest{
		Index: []string{indexName(req.QueryId)},
		Body: strings.NewReader(fmt.Sprintf(`{
			"size": 10000,
			"query": {
        "bool": {
          "must": [
            { "term": { "StepId": "%s" } }
          ]
        }
      }
		}`, req.StepId)),
	}

	res, err := searchReq.Do(context.Background(), r.esclient)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting response: %s", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing the response body: %s", err))
		} else {
			// Print the response status and error information.
			return nil, errors.New(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			))
		}
	}


	// Print the response status, number of results, and request duration.
	r2 := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&r2); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r2["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r2["took"].(float64)),
	)

	// Print the ID and document source for each hit.
	results := make([]string, 0)
	for _, hit := range r2["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])

		source := hit.(map[string]interface{})["_source"]

		record := resultRecord{
			StepId: source.(map[string]interface{})["StepId"].(string),
			OtherId: source.(map[string]interface{})["OtherId"].(string),
			RecordId: source.(map[string]interface{})["RecordId"].(string),
		}
		if record.OtherId != "" {
			results = append(results, record.OtherId)
		} else {
			results = append(results, record.RecordId)
		}
	}

	return results, nil
}

func indexName(queryId string) string {
	return fmt.Sprintf("asearch_results_%s", queryId)
}
