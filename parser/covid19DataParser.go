package parser

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ritwickdey/covid-19-india-golang/model"
)

type Covid19DataParser interface {
	DownloadAndParse(url string) (model.Covid19StatMap, error)
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

	fmt.Println("Total tr: ", trLen)

	for i := 0; i < trLen; i++ {
		c.processTRSelection(i, tr)

		//I know it's BAD. but I've no other workaround. :| ...
		if strings.ToLower(strings.Replace(c.currentParsedData.StateName, " ", "", -1)) == "westbengal" {
			tr = tr.Next()

			isReassigned := strings.Contains(tr.Find("td:nth-child(2)").Text(), "reassigned")
			if isReassigned {
				other := c.processTRSelection(i, tr)

				delete(c.data, other.StateName)
				other.StateName = "Other"
				c.data["Other"] = other

			}

			break
		}

		tr = tr.Next()
	}

	return c.data, err
}

func (c *covid19DataParser) processTRSelection(index int, selection *goquery.Selection) model.Covid19Stat {
	c.currentParsedData = model.Covid19Stat{}
	selection.Find("td").Each(c.processTDSelection)

	c.data[c.currentParsedData.StateName] = c.currentParsedData

	return c.currentParsedData

}

func (c *covid19DataParser) processTDSelection(index int, selection *goquery.Selection) {

	re := regexp.MustCompile(`\*|\#`)
	text := re.ReplaceAllString(selection.Text(), "")

	switch index {
	case 0:
		break
	case 1:
		c.currentParsedData.StateName = text
		break
	case 2:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.ActiveCase = i
		break
	case 3:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.Cured = i
		break
	case 4:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.Death = i
	case 5:
		i, _ := strconv.Atoi(text)
		c.currentParsedData.ConfirmedCase = i
		break
	}
}
