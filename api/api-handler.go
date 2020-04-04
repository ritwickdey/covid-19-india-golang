package api

import "github.com/ritwickdey/covid-19-india-go-lang/model"

type Service interface {
	FetchAllData() (model.Covid19StatMapDateWise, error)
	FetchByDate(date string) (model.Covid19StatMap, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s *service) FetchAllData() (model.Covid19StatMapDateWise, error) {
	return model.DataCache.GetCache(), nil
}

func (s *service) FetchByDate(date string) (model.Covid19StatMap, error) {
	return model.DataCache.GetCache()[date], nil
}
