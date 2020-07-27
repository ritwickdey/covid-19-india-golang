package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ritwickdey/covid-19-india-golang/api"
	"github.com/ritwickdey/covid-19-india-golang/model"
	"github.com/ritwickdey/covid-19-india-golang/parser"
)

var WEB_END_POINT = "https://www.mohfw.gov.in"
var FILE_PATH = "./output-stats.json"

func main() {

	args := os.Args[1:]

	if len(args) > 0 {
		FILE_PATH = args[0]
	}

	fmt.Println(FILE_PATH)

	existingData, err := readExistingData()
	throwIfErr(err)
	model.DataCache.UpdateCache(existingData)
	go fetchDataPeriodically()

	service := api.NewService()
	mux := CORS(api.MakeHTTPHandler(service))

	serverAddress := ":5566"

	log.Println("Server started with", serverAddress)
	log.Fatalln(http.ListenAndServe(serverAddress, mux))

}

func CORS(h http.Handler) http.Handler {

	sites := make(map[string]string)
	sites["https://novelcoronaindia.info"] = "https://novelcoronaindia.info"
	sites["http://localhost:3000"] = "http://localhost:3000"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if sites[origin] != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		}

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
	loc, _ := time.LoadLocation("Asia/Kolkata")
	todayKey := time.Now().In(loc).Format("02-01-2006")
	p := parser.NewCovid19APIParser()
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
