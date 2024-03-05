package multi_source_game_scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

type GameData struct {
	url, image, title, releaseDate string
	price, userScore               float32
	metascore                      int
	tags                           []string
}

func ScrapeMetacritic() []GameData {

	var metacriticGames []GameData

	collector := colly.NewCollector(
		colly.MaxDepth(1),
	)

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting:", request.URL)
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("Page visited:", response.Request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnScraped(func(request *colly.Response) {
		fmt.Println(request.Request.URL, "scraped!")
	})

	// Goes into link of each element on page and scrapes the game data from their specific site
	collector.OnHTML("a.c-finderProductCard_container", func(element *colly.HTMLElement) {
		metacriticGames = append(metacriticGames, ScrapeSingleMetacriticGame("https://www.metacritic.com"+element.Attr("href")))
	})

	err := collector.Visit("https://www.metacritic.com/browse/game/")
	if err != nil {
		log.Println("Error when visiting site; pointer set to", err)
		return metacriticGames
	}
	return metacriticGames
}

func ScrapeSingleMetacriticGame(link string) GameData {

	var gameData GameData
	gameData.url = link
	log.Println(gameData.url)

	collector := colly.NewCollector()

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting:", request.URL)
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Println("Page visited:", response.Request.URL)
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnScraped(func(request *colly.Response) {
		fmt.Println(request.Request.URL, "scraped!")
	})

	collector.OnHTML("div.c-productHero_title", func(element *colly.HTMLElement) {
		gameData.title = element.ChildText("div")
		log.Println(gameData.title)

	})

	return gameData
}
