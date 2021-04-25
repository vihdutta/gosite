package quote

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
func QuoteGen() string {
	rand.Seed(time.Now().Unix())
	jsonFile, _ := os.Open("quotes.json")
	bodyBytes, _ := ioutil.ReadAll(jsonFile)

	var quotes []RandQuote
	json.Unmarshal(bodyBytes, &quotes)

	return ParseQuotes(quotes)
}

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
