package transport

import (
	"context"

	"pharmacy/app/implementation"
	"pharmacy/app/medicine"

	"github.com/go-kit/kit/endpoint"
)

type CreateRequest struct {
	Medicine medicine.Medicine `json:"medicine"`
}
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}
type GetByIDResponse struct {
	Medicine medicine.Medicine `json:"medicine"`
	Err      error             `json:"-"`
}

type SellRequest struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}
type SellResponse struct {
	OK  bool  `json:"ok"`
	Err error `json:"-"`
}

type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	Sell    endpoint.Endpoint
}

func MakeEndpoints(s implementation.Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
		Sell:    makeSellEndpoint(s),
	}
}

func makeCreateEndpoint(s implementation.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Medicine)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s implementation.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		m, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Medicine: m, Err: err}, nil
	}
}

func makeSellEndpoint(s implementation.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SellRequest)
		err := s.Sell(ctx, req.ID, req.Count)
		return SellResponse{OK: err == nil, Err: err}, nil
	}
}
