package multi_source_game_scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

type GameData struct {
	url, image, name, releaseDate string
	price, userScore              float32
	metascore                     int
	tags                          []string
}

func ScrapeMetacritic() []GameData {

	var metacriticGames []GameData

	collector := colly.NewCollector()

	err := collector.Visit("https://scrapeme.live/shop/")
	if err != nil {
		log.Println("Error when visiting site; pointer set to ", err)
		return metacriticGames
	}

	log.Println("Collector visited website successfully!")

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting: ", request.URL)
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("Page visited: ", response.Request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	collector.OnScraped(func(request *colly.Response) {
		fmt.Println(request.Request.URL, " scraped!")
	})

	collector.OnHTML("li.product", func(element *colly.HTMLElement) {
		game := GameData{}
		//game.url = element.Attr("href")
		game.url = element.ChildAttr("a", "href")
		fmt.Println(game.url)

		metacriticGames = append(metacriticGames, game)
	})

	return metacriticGames
}
