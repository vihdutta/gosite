package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type RandQuote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

/*
func main() {
	fmt.Println(QuoteGen())
}
*/
func QuoteGen() RandQuote {
	rand.Seed(time.Now().Unix())

	jsonFile, err := os.Open("static/json/quotes.json")

	if err != nil {
		fmt.Println(err)
	}

	bodyBytes, _ := ioutil.ReadAll(jsonFile)

	var quotes []RandQuote
	json.Unmarshal(bodyBytes, &quotes)
	quote := quotes[rand.Intn(len(quotes))]

	if quote.Author == "" {
		quote.Author = "Anonymous"
	}
	return quote
}

/*
func ParseQuotes(quotes []RandQuote) string {
	quote := quotes[rand.Intn(len(quotes))]

	if quote.Author == "" {
		quote.Author = "Anonymous"
	}

	fquote := fmt.Sprintf("\"%+v\" - %+v", quote.Text, quote.Author)
	fquotelen := len(fquote)

	switch {
	case fquotelen > 80:
		return ParseQuotes(quotes)
	default:
	}
	return fquote
}
*/
