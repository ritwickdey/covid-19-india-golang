package parser

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/ritwickdey/covid-19-india-golang/model"
)

type Covid19DataParser interface {
	DownloadAndParse(url string) (model.Covid19StatMap, error)
	ParseFromReader(r io.Reader) (model.Covid19StatMap, error)
}

func NewCovid19DataParser() Covid19DataParser {
	return &covid19DataParser{
		data:              model.Covid19StatMap{},
		currentParsedData: model.Covid19Stat{},
	}
}

type covid19DataParser struct {
	data              model.Covid19StatMap
	currentParsedData model.Covid19Stat
}

func (c *covid19DataParser) DownloadAndParse(url string) (model.Covid19StatMap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return c.data, err
	}

	return c.ParseFromReader(resp.Body)
}

func (c *covid19DataParser) ParseFromReader(r io.Reader) (model.Covid19StatMap, error) {
	html, err := goquery.NewDocumentFromReader(r)
	tableHtml := html.Find("#state-data .data-table > table > tbody")

	tr := tableHtml.Find("tr")
	trLen := tr.Length()

	fmt.Println(trLen)

	for i := 0; i < 30; i++ { // 30 row data
		c.processTRSelection(i, tr)
		tr = tr.Next()
	}

	return c.data, err
}

func (c *covid19DataParser) processTRSelection(index int, selection *goquery.Selection) {
	c.currentParsedData = model.Covid19Stat{}
	selection.Find("td").Each(c.processTDSelection)

	c.data[c.currentParsedData.StateName] = c.currentParsedData

}
func (c *covid19DataParser) processTDSelection(index int, selection *goquery.Selection) {
	text := selection.Text()
	switch index {
	case 0:
		break
	case 1:
		c.currentParsedData.StateName = text
		break
	case 2:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.ConfirmedCase = i
		break
	case 3:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.Cured = i
		break
	case 4:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.Death = i
		break
	}
}
