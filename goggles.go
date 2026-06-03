package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

const defaultUserAgent = "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/131.0.0.0 Safari/537.36"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("goggles", flag.ContinueOnError)
	selector := fs.String("selector", "body", "CSS selector to extract")
	timeout := fs.Duration("timeout", 10*time.Second, "request timeout")
	userAgent := fs.String("ua", defaultUserAgent, "User-Agent header")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("please pass a valid URL as the first argument")
	}
	url := fs.Arg(0)

	client := http.Client{Timeout: *timeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", *userAgent)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("non-success status: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var iterErr error
	doc.Find(*selector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		html, err := s.Html()
		if err != nil {
			iterErr = err
			return false
		}
		text, err := html2text.FromString(html, html2text.Options{PrettyTables: true})
		if err != nil {
			iterErr = err
			return false
		}
		fmt.Fprintln(out, text)
		return true
	})
	return iterErr
}
