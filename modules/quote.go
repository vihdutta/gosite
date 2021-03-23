package quote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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
	req, err := http.Get("https://type.fit/api/quotes")

	if err != nil {
		fmt.Print(err.Error())
	}

	defer req.Body.Close()
	bodyBytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

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
