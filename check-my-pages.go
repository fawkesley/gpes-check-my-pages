package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/asciimoo/colly"
)

func main() {

	urlsFileName := flag.String("urls", "urls.csv", "Name of the csv file with the urs in the first column")
	isHTTP := flag.Bool("http", false, "Http response codes")
	isAnalytics := flag.Bool("analytics", false, "Correct analytics tag in the html")
	isCanonical := flag.Bool("canonical", false, "Canonical URLS in the ")
	isClear := flag.Bool("clear", false, "Remove files created by this script")
	flag.Parse()

	allUrlsCsv := readCsvFile(*urlsFileName)

	allUrls := csvFirstColumnToSlice(allUrlsCsv)

	c := colly.NewCollector()
	// c.AllowedDomains = []string{"localhost", "greenpeace.es", "archivo.greenpeace.es"}

	if *isHTTP == true {

		httpResponses, httpErr := os.OpenFile("httpResponses.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if httpErr != nil {
			panic(httpErr)
		}
		defer httpResponses.Close()

		c.OnResponse(func(r *colly.Response) {
			lineResponse := fmt.Sprintf("%s,%v\n", r.Request.URL.String(), r.StatusCode)
			if _, err := httpResponses.WriteString(lineResponse); err != nil {
				panic(err)
			}

		})
	}

	if *isAnalytics == true {

		analytics, analyticsErr := os.OpenFile("analytics.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if analyticsErr != nil {
			panic(analyticsErr)
		}
		defer analytics.Close()

		c.OnResponse(func(r *colly.Response) {
			body := string(r.Body)
			foundUA := searchInString(body, `UA-\d{5,8}-\d{1,2}`)
			lineResponse := fmt.Sprintf("%s,%s\n", r.Request.URL.String(), foundUA)
			if _, err := analytics.WriteString(lineResponse); err != nil {
				panic(err)
			}
		})
	}

	if *isCanonical == true {

		canonical, canonicalErr := os.OpenFile("canonicals.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if canonicalErr != nil {
			panic(canonicalErr)
		}
		defer canonical.Close()

		c.OnHTML("link[rel=canonical]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			lineCanonical := fmt.Sprintf("%s,%s\n", e.Request.URL.String(), link)
			if _, err := canonical.WriteString(lineCanonical); err != nil {
				panic(err)
			}
		})
	}

	if *isClear == true {

		os.Remove("httpResponses.csv")
		os.Remove("analytics.csv")
		os.Remove("canonicals.csv")
	}

	// Open URLs file
	for _, v := range allUrls {
		c.Visit(v)
		time.Sleep(time.Millisecond * 100)
	}

}