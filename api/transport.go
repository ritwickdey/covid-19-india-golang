package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeGetAllStatsEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return s.FetchAllData()
	}
}

func MakeGetStatsByDateEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		var date = request.(string)
		return s.FetchByDate(date)
	}
}

func DecodeGetAllDataReq(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeGetStatsByDateReq(_ context.Context, r *http.Request) (interface{}, error) {

	date := mux.Vars(r)["date"]
	if date == "" {
		return nil, errors.New("Date is missing")
	}

	return date, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

func MakeHTTPHandler(s Service) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{}

	r.Methods("GET").Path("/covid19/all").Handler(
		httptransport.NewServer(MakeGetAllStatsEndpoint(s),
			DecodeGetAllDataReq,
			EncodeResponse,
			options...,
		),
	)
	r.Methods("GET").Path("/covid19/date/{date}").Handler(
		httptransport.NewServer(MakeGetStatsByDateEndpoint(s),
			DecodeGetStatsByDateReq,
			EncodeResponse,
			options...,
		),
	)

	return r
}
