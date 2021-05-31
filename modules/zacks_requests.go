package modules

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PuerkitoBio/goquery"
)

var stocks []string
var stockValues []string

func ZacksRequests(tempFile string) {
	f, err := excelize.OpenFile(tempFile)
	if err != nil {
		fmt.Println(err)
	}

	name := f.GetSheetMap()

	rows, err := f.GetRows(name[1])
	if err != nil {
		fmt.Println(err)
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		stocks = append(stocks, row[0])
	}

	for _, stock := range stocks {
		url := "https://www.zacks.com/stock/quote/" + stock + "?q=" + stock

		r, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}

		rdoc, err := goquery.NewDocumentFromReader(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		rdoc.Find("#quote_ribbon_v2 > div.quote_rank_summary > div:nth-child(1) > p").Each(func(i int, s *goquery.Selection) {
			stockValue := strings.TrimSpace(strings.Split(s.Text(), "of")[0])
			statement := stock + ": " + stockValue
			stockValues = append(stockValues, statement)
		})
	}

	analysisFile, err := os.OpenFile("analysis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer analysisFile.Close()

	for _, statement := range stockValues {
		if _, err = analysisFile.WriteString(statement + "\n"); err != nil {
			fmt.Println(err)
		}
	}
}
