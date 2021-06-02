package modules

import (
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gocolly/colly/v2"
)

type Statement struct {
	stock       string
	stockRating string
	stockStats  string
}

func ZacksRequests(tempFile, formFileName string) {
	var stocks []string
	var stockStatements []string
	var unfoundItems []string

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
		found := false

		url := "https://www.zacks.com/stock/quote/" + stock + "?q=" + stock

		collector := colly.NewCollector(
			colly.AllowedDomains("zacks.com", "www.zacks.com"),
		)

		stockData := Statement{}

		//stock rating
		collector.OnXML("//*[@id='quote_ribbon_v2']/div[2]/div[1]/p/text()", func(element *colly.XMLElement) {
			if strings.TrimSpace(element.Text) != "" {
				stockData.stock = stock
				stockData.stockRating = strings.TrimSpace(element.Text)
				found = true
			}
		})

		//stock value, growth, momentum, vgm
		collector.OnXML("//*[@id='quote_ribbon_v2']/div[2]/div[2]/p", func(element *colly.XMLElement) {
			replacer := strings.NewReplacer("Value", "", "Growth", "", "Momentum", "", "VGM", "", "|", "")
			stockStats := strings.TrimSpace(replacer.Replace(element.Text))
			stockStats = strings.ReplaceAll(stockStats, "  ", "")
			stockData.stockStats = stockStats
		})

		collector.Visit(url)

		statement := stockData.stock + ": " + stockData.stockRating + " " + stockData.stockStats

		if strings.TrimSpace(stockData.stock) != "" {
			stockStatements = append(stockStatements, statement)
		}

		if !found {
			unfoundItems = append(unfoundItems, stock)
		}
	}

	analysisFile, err := os.OpenFile(formFileName+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	for _, statement := range stockStatements {
		if _, err = analysisFile.WriteString(statement + "\n"); err != nil {
			fmt.Println(err)
		}
	}

	analysisFile.WriteString(fmt.Sprintf("Unfound items: %s", strings.Join(unfoundItems, ", ")))

	defer analysisFile.Close()
}
