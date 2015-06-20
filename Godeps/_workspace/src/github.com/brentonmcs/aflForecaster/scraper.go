package AflForecaster

import (
	"aflForecasterRunner/Godeps/_workspace/src/github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
)

func ScrapePages() {

	activeRounds := getActiveRounds()

	baseUri := "http://footyforecaster.com/AFL/RoundForecast/%d_Round_%d"
	if len(activeRounds) == 0 {
		year, round := getCurrentRound()
		scrapeResults(fmt.Sprintf(baseUri, year, round+1), round+1, year)
	}

	for _, aR := range activeRounds {

		scrapeResults(fmt.Sprintf(baseUri, aR.year, aR.round), aR.round, aR.year)

		fmt.Printf("Finished scraping year: %d round :%d", aR.year, aR.round)
	}
}

func goToPage(uri string) *goquery.Document {
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func scrapeResults(uri string, round, year int) {

	doc := goToPage(uri)

	doc.Find(".details").Each(func(i int, s *goquery.Selection) {

		tableBody := s.Find("tbody")
		forecast := tableBody.Find("tr:nth-child(3) td:nth-child(2) td:first-child").Text()
		percentage := tableBody.Find("tr:nth-child(1) > td:nth-child(2)").Text()
		result := tableBody.Find("tr:nth-child(4) > td:nth-child(2) ").Text()

		var forecastModel = parseForecast(forecast, percentage, result, round, year)
		AddForecast(&forecastModel)
	})
}
