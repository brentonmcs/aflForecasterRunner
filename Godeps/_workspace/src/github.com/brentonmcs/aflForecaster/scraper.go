package AflForecaster

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ScrapePages() {

	activeRounds := getActiveRounds()

	baseUri := "http://footyforecaster.com/AFL/RoundForecast/%d_Round_%d"
	if len(activeRounds) == 0 {
		year, round := getCurrentRound()
		scrapeResults(fmt.Sprintf(baseUri, year, round +1), round +1, year)
	}


	scapeSportsBet()

	for _, aR := range activeRounds {

		scrapeResults(fmt.Sprintf(baseUri, aR.year, aR.round), aR.round, aR.year)

		fmt.Printf("Finished scraping year: %d round :%d", aR.year, aR.round)
	}

}

func scapeSportsBet() {

	clearPrices()
	doc := goToPage("http://www.sportsbet.com.au/betting/australian-rules?QuickLinks")

	headerAndAccord := doc.Find(".accordion-main").Find(".bettypes-header, .accordion-body")
	var err error
	var currentDate time.Time

	headerAndAccord.Each(func(i int, s *goquery.Selection) {
		if s.HasClass("bettypes-header") {
			date := s.Find(".date").Text()
			currentDate, err = time.Parse("02/01/2006", strings.Split(date, " ")[1])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			buttons := s.Find(".market-buttons")

			homeTeam := parsePrices(buttons.First())
			awayTeam := parsePrices(buttons.Last())

			matchPrice := MatchPrices{HomeTeam:homeTeam, AwayTeam: awayTeam, MatchDate: currentDate}
			fmt.Println(matchPrice)
			addPriceRecord(&matchPrice)
		}
	});


}

func parsePrices(prices *goquery.Selection) PriceModel {

	priceBox := prices.Find(".price-link")

	headToHeadBox := priceBox.First()
	headToHeadPrice := headToHeadBox.Find(".odd-val").Text()

	priceBox = priceBox.Next()
	under40Price := priceBox.First().Find(".odd-val").Text()

	priceBox = priceBox.Next()
	over40Price := priceBox.First().Find(".odd-val").Text()

	return PriceModel{Name:headToHeadBox.Find(".team-name").Text(),
		HeadToHead: strToPrice(headToHeadPrice),
		Under39: strToPrice(under40Price),
		Over40: strToPrice(over40Price)}
}

func strToPrice(strPrice string) float32 {
	price, err := strconv.ParseFloat(strings.Trim(strPrice, "\n"), 32)
	if err != nil {
		log.Fatal(err)
	}
	return float32(price)
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


