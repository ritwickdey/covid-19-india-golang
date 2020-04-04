package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ritwickdey/covid-19-india-go-lang/api"
	"github.com/ritwickdey/covid-19-india-go-lang/model"
	"github.com/ritwickdey/covid-19-india-go-lang/parser"
)

var WEB_END_POINT = "https://www.mohfw.gov.in"
var FILE_PATH = "./output-stats.json"

func main() {
	exitingData, err := readExistingData()
	throwIfErr(err)
	model.DataCache.UpdateCache(exitingData)
	go fetchDataPeriodically()

	service := api.NewService()
	mux := CORS(api.MakeHTTPHandler(service))

	serverAddress := ":5566"

	log.Println("Server started with", serverAddress)
	log.Fatalln(http.ListenAndServe(serverAddress, mux))

}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func fetchDataPeriodically() {
	model.DataCache.UpdateCache(dataParserFromOfficialSite())

	for range time.NewTicker(30 * time.Minute).C {
		model.DataCache.UpdateCache(dataParserFromOfficialSite())
	}
}

func dataParserFromOfficialSite() model.Covid19StatMapDateWise {
	todayKey := time.Now().Format("02-01-2006")
	p := parser.NewCovid19DataParser()
	currentData, err := p.DownloadAndParse(WEB_END_POINT)
	throwIfErr(err)

	existingData, err := readExistingData()
	throwIfErr(err)

	existingData[todayKey] = currentData

	optJson, err := json.Marshal(existingData)
	throwIfErr(err)

	err = ioutil.WriteFile(FILE_PATH, optJson, 0644)
	throwIfErr(err)

	log.Println("data fetched from official site")

	return existingData
}

func readExistingData() (model.Covid19StatMapDateWise, error) {
	dataBytes, err := ioutil.ReadFile(FILE_PATH)
	if err != nil {
		dataBytes = []byte("{}")
	}

	output := model.Covid19StatMapDateWise{}

	err = json.Unmarshal(dataBytes, &output)

	return output, err

}

func throwIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
