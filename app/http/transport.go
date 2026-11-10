package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"pharmacy/app/medicine"
	"pharmacy/app/transport"

	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(e transport.Endpoints) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /medicines", kithttp.NewServer(
		e.Create, decodeCreate, encodeCreate, kithttp.ServerErrorEncoder(errorEncoder),
	))
	mux.Handle("GET /medicines/{id}", kithttp.NewServer(
		e.GetByID, decodeGetByID, encodeGetByID, kithttp.ServerErrorEncoder(errorEncoder),
	))
	mux.Handle("POST /medicines/{id}/sell", kithttp.NewServer(
		e.Sell, decodeSell, encodeSell, kithttp.ServerErrorEncoder(errorEncoder),
	))
	return mux
}

func decodeCreate(_ context.Context, r *http.Request) (interface{}, error) {
	var m medicine.Medicine
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}
	return transport.CreateRequest{Medicine: m}, nil
}

func decodeGetByID(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetByIDRequest{ID: r.PathValue("id")}, nil
}

func decodeSell(_ context.Context, r *http.Request) (interface{}, error) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		var body struct {
			Count int `json:"count"`
		}
		if decErr := json.NewDecoder(r.Body).Decode(&body); decErr != nil {
			return nil, err
		}
		count = body.Count
	}
	return transport.SellRequest{ID: r.PathValue("id"), Count: count}, nil
}

func encodeCreate(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(transport.CreateResponse)
	if resp.Err != nil {
		errorEncoder(context.Background(), resp.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(map[string]string{"id": resp.ID})
}

func encodeGetByID(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(transport.GetByIDResponse)
	if resp.Err != nil {
		errorEncoder(context.Background(), resp.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp.Medicine)
}

func encodeSell(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(transport.SellResponse)
	if resp.Err != nil {
		errorEncoder(context.Background(), resp.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]bool{"ok": resp.OK})
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, medicine.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, medicine.ErrInvalidMedicine), errors.Is(err, medicine.ErrInvalidCount):
		status = http.StatusBadRequest
	case errors.Is(err, medicine.ErrOutOfStock):
		status = http.StatusConflict
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
