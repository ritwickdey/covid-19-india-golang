package api

import (
	"time"

	"github.com/ritwickdey/covid-19-india-golang/model"
)

type Service interface {
	FetchAllData() (model.Covid19StatMapDateWise, error)
	FetchByDate(date string) (model.Covid19StatMap, error)
	FetchByDateRange(from string, to string) (model.Covid19StatMapDateWise, error)
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

func (s *service) FetchByDateRange(fromStr string, toStr string) (model.Covid19StatMapDateWise, error) {
	result := model.Covid19StatMapDateWise{}
	from, err := time.Parse("02-01-2006", fromStr)

	if err != nil {
		return result, err
	}

	to, err := time.Parse("02-01-2006", toStr)

	if err != nil {
		return result, err
	}

	for i := from; to.Sub(i).Hours() >= 0; i = i.Add(time.Hour * 24) {
		d := i.Format("02-01-2006")
		covid19StatMap, err := s.FetchByDate(d)
		if err != nil {
			return result, err
		}
		result[d] = covid19StatMap
	}

	return result, nil

}
