package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Please pass a valid URL as first argument.")
	}
	url := os.Args[1]
	var selector = "body"
	if len(os.Args) > 2 {
		selector = os.Args[2]
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		html, _ := s.Html()
		text, err := html2text.FromString(string(html), html2text.Options{PrettyTables: true})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(text)
	})
}
