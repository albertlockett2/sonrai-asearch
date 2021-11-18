package resultsdao

import (
	"github.com/elastic/go-elasticsearch/v7"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
)

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
	// TODO
	return nil
}


func (r *ResultsDao) WriteResults(results *gen.ResultRecord) error {
	// TODO
	return nil
}