package api

import (
	"time"

	"github.com/ritwickdey/covid-19-india-golang/model"
)

type Service interface {
	FetchAllData() (model.Covid19StatMapDateWise, error)
	FetchByDate(date string) (model.Covid19StatMap, error)
	FetchByDateRange(from string, to string) (model.Covid19StatMapDateWise, error)
	FetchByDateRangeFormated(from time.Time, to time.Time) (model.FormatedStatResult, error)
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
	from, err := time.Parse(model.DateFormatPattern, fromStr)

	if err != nil {
		return result, err
	}

	to, err := time.Parse(model.DateFormatPattern, toStr)

	if err != nil {
		return result, err
	}

	for i := from; to.Sub(i).Hours() >= 0; i = i.Add(time.Hour * 24) {
		d := i.Format(model.DateFormatPattern)
		covid19StatMap, err := s.FetchByDate(d)
		if err != nil {
			return result, err
		}
		result[d] = covid19StatMap
	}

	return result, nil

}

func (s *service) FetchByDateRangeFormated(from time.Time, to time.Time) (model.FormatedStatResult, error) {
	result := model.FormatedStatResult{}

	for i := from; to.Sub(i).Hours() >= 0; i = i.Add(time.Hour * 24) {
		d := i.Format(model.DateFormatPattern)
		covid19StatMap, err := s.FetchByDate(d)
		if err != nil {
			return result, err
		}

		r := model.FormatedStatData{
			Date:      d,
			StateWise: []model.StateWiseFormatedStatData{},
		}

		for _, stateData := range covid19StatMap {
			r.Confirmed += stateData.ConfirmedCase
			r.Death += stateData.Death
			r.Recovered += stateData.Cured
			r.StateWise = append(r.StateWise, model.StateWiseFormatedStatData{
				StateName: stateData.StateName,
				Confirmed: stateData.ConfirmedCase,
				Recovered: stateData.Cured,
				Death:     stateData.Death,
				Active:    stateData.ConfirmedCase - stateData.Cured - stateData.Death,
			})
		}
		r.Active = r.Confirmed - r.Recovered - r.Death
		result.Data = append(result.Data, r)
	}

	return result, nil

}
