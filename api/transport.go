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

func MakeGetStatsByDateRangeEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		var dates = request.([]string)
		return s.FetchByDateRange(dates[0], dates[1])
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

func DecodeGetStatsByDateRangeReq(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	startDate := vars["startDate"]
	endDate := vars["endDate"]
	if startDate == "" {
		return nil, errors.New("start date is missing")
	}
	if endDate == "" {
		return nil, errors.New("end date is missing")
	}
	dates := []string{startDate, endDate}
	return dates, nil
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

	r.Methods("GET").Path("/covid19/dateRange/{startDate}/{endDate}").Handler(
		httptransport.NewServer(MakeGetStatsByDateRangeEndpoint(s),
			DecodeGetStatsByDateRangeReq,
			EncodeResponse,
			options...,
		),
	)

	return r
}
