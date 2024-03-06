package multi_source_game_scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type GameData struct {
	url, image, title, releaseDate, description, rating, developer, publisher string
	userScore                                                                 float32
	metascore                                                                 int
	tags, platforms                                                           []string
}

func ScrapeMetacriticPage(link string) []GameData {

	var metacriticGames []GameData

	collector := colly.NewCollector(
	//colly.Async(true),
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
		game := ScrapeSingleMetacriticGame("https://www.metacritic.com" + element.Attr("href"))
		metacriticGames = append(metacriticGames, game)
	})

	//if HasNextPage(collector) {
	//	metacriticGames = extend(metacriticGames, ScrapeMetacriticPage(ReturnMetacriticNextPage(link, collector)))
	//}

	err := collector.Visit("https://www.metacritic.com/browse/game/")
	if err != nil {
		log.Println("Error when visiting site; pointer set to", err)
		return metacriticGames
	}
	return metacriticGames
}

// Only returns User Score, Review Score, Title, Publisher, and Release Date
func ScrapeSingleMetacriticGame(link string) GameData {

	var gameData GameData
	gameData.url = link

	collector := colly.NewCollector(
		colly.Async(true),
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

	//Title
	collector.OnHTML("div.c-productHero_title", func(element *colly.HTMLElement) {
		gameData.title = element.ChildText("div")
		log.Println(gameData.title)
	})

	//Release Date
	collector.OnHTML("div.g-text-xsmall", func(element *colly.HTMLElement) {
		gameData.releaseDate = element.ChildText("span.u-text-uppercase")
		log.Println(gameData.releaseDate)
	})

	//Metacritic and User Scores
	collector.OnHTML("div.c-productScoreInfo_scoreNumber", func(element *colly.HTMLElement) {

		//if has period then must be Float and must be user score
		if strings.Contains(element.ChildText("span"), ".") {

			float64Temp, float64ConversionError := strconv.ParseFloat(element.ChildText("span"), 32)

			if float64ConversionError != nil {
				log.Println("Error when converting String to Float:", float64ConversionError)
			}

			gameData.userScore = float32(float64Temp)

			log.Println(gameData.userScore)
		} else {
			gameData.metascore, _ = strconv.Atoi(element.ChildText("span"))
			log.Println(gameData.metascore)
		}

	})

	//Rating and Developer
	collector.OnHTML("div.c-heroMetadata", func(element *colly.HTMLElement) {

		if len(gameData.rating) < 1 {
			devAndRating := strings.Split(element.ChildText("span"), "\n")
			gameData.rating = devAndRating[0]
			gameData.publisher = strings.TrimSpace(devAndRating[2])
			log.Println(gameData.rating)
			log.Println(gameData.publisher)
		}

	})

	//Platforms
	collector.OnHTML("div.c-gamePlatformLogo", func(element *colly.HTMLElement) {
		// if not already in gameData.platforms then add data
		if !(strings.Contains(strings.Join(gameData.platforms, ","), element.ChildText("div.g-text-medium"))) {
			gameData.platforms = append(gameData.platforms, element.ChildText("div.g-text-medium"))
			log.Println(gameData.platforms)
		} else if !(strings.Contains(strings.Join(gameData.platforms, ","), element.ChildText("title"))) {
			gameData.platforms = append(gameData.platforms, element.ChildText("title"))
			log.Println(gameData.platforms)
		}
	})

	err := collector.Visit(link)
	if err != nil {
		log.Println("Error when visiting site; pointer set to", err)
		return gameData
	}

	return gameData
}

// ReturnMetacriticNextPage returns blank if page not found
func ReturnMetacriticNextPage(link string, collector *colly.Collector) string {

	if HasNextPage(collector) {
		if strings.HasSuffix(link, "page=%d") {
			tempSubSplice := strings.SplitN(link, "page=", 2)
			link = tempSubSplice[0]

			pageNumber, err := strconv.Atoi(tempSubSplice[1])

			if err != nil {
				log.Println("Cannot convert string to integer:", err)
			}
			link = link + "page=" + strconv.Itoa(pageNumber+1)
		} else {
			//if not has page suffix and has no next page then must be first page
			link = link + "?page=2"
		}
	}

	return ""

}

func HasNextPage(collector *colly.Collector) bool {

	var nextPageExists = false

	collector.OnHTML("c-navigationPagination_item", func(element *colly.HTMLElement) {
		if !(element.Attr("c-navigationPagination_item--next") == "disabled") {
			nextPageExists = true
		}
	})
	return nextPageExists
}
