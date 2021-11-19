package main

import (
	"context"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"github.com/sonraisecurity/sonrai-asearch/src/results/dao"
)

type ResultsService struct {
	gen.UnimplementedResultsServer
	dao *dao.ResultsDao
}

func NewResultsService() (*ResultsService, error) {
	r, err := dao.NewResultsDao()
	if err != nil {
		return nil, err
	}

	return &ResultsService{
		dao: r,
	}, nil
}

func (r *ResultsService) GetResults(ctx context.Context, req *gen.ResultsRequest) (*gen.ResultsResponse, error) {
	results, err := r.dao.GetResults(req)
	if err != nil {
		return nil, err
	}
	return &gen.ResultsResponse{
		Ids: results,
	}, nil
}
