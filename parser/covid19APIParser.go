package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ritwickdey/covid-19-india-golang/model"
)

func NewCovid19APIParser() Covid19DataParser {
	return &covid19APIParser{
		data: model.Covid19StatMap{},
	}
}

type covid19APIParser struct {
	data model.Covid19StatMap
}

func (c *covid19APIParser) DownloadAndParse(_ string) (model.Covid19StatMap, error) {
	res, err := http.Get("https://www.mohfw.gov.in/data/datanew.json")
	if err != nil {
		return c.data, err
	}

	defer res.Body.Close()

	type ApiRes struct {
		StateName string `json:"state_name"`
		Active    string `json:"new_active"`
		Positive  string `json:"new_positive"`
		Cured     string `json:"new_cured"`
		Death     string `json:"new_death"`
	}

	rawData := make([]ApiRes, 0)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&rawData)

	for _, apiData := range rawData {
		active, _ := strconv.Atoi(apiData.Active)
		cured, _ := strconv.Atoi(apiData.Cured)
		death, _ := strconv.Atoi(apiData.Death)
		confirmed, _ := strconv.Atoi(apiData.Positive)
		stateName := removeStars(apiData.StateName)
		fmt.Println(stateName)

		if stateName == "" {
			continue
		}

		if stateName == "Telengana" {
			stateName = "Telangana"
		}

		c.data[stateName] = model.Covid19Stat{
			StateName:     stateName,
			ActiveCase:    active,
			Cured:         cured,
			Death:         death,
			ConfirmedCase: confirmed,
		}
	}

	return c.data, err

}

func removeStars(str string) string {
	index := strings.Index(str, "*")
	if index == -1 {
		return str
	}

	return str[:index]
}
