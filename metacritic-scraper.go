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

func ScrapeMetacriticPage(currentLink string) []GameData {

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

	// Goes into currentLink of each element on page and scrapes the game data from their specific site
	collector.OnHTML("a.c-finderProductCard_container", func(element *colly.HTMLElement) {
		game := ScrapeSingleMetacriticGame("https://www.metacritic.com" + element.Attr("href"))
		metacriticGames = append(metacriticGames, game)
	})

	//Goes to next page
	collector.OnHTML("span.c-navigationPagination_itemButtonContent", func(element *colly.HTMLElement) {

		//#arrowChevron means there is a next page
		//if not next page then this element should be #arrowChevronDisabled
		if element.ChildAttr("use", "href") == "#arrowChevron" {
			log.Println(element.ChildAttr("use", "href"))
			tempSubSplice := strings.SplitN(currentLink, "page=", 2)

			currentPageNumber, err := strconv.Atoi(tempSubSplice[1])

			if err != nil {
				log.Println("Error when visiting site; pointer set to", err)
			}

			nextPageNumber := currentPageNumber + 1

			ScrapeMetacriticPage(tempSubSplice[0] + "page=" + strconv.Itoa(nextPageNumber))

		}

	})

	err := collector.Visit(currentLink)
	if err != nil {
		log.Println("Error when visiting site; pointer set to", err)
		return metacriticGames
	}
	return metacriticGames
}

// Only returns User Score, Review Score, Title, Publisher, and Release Date
func ScrapeSingleMetacriticGame(currentLink string) GameData {

	var gameData GameData
	gameData.url = currentLink

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

		//if gameData.rating is empty then get game rating
		if len(gameData.rating) < 1 {
			devAndRating := strings.Split(element.ChildText("span"), "\n")
			gameData.rating = devAndRating[0]

			// if devAndRating is greater than 2 then there is a publisher record assigned to devAndRating[2]
			if len(devAndRating) > 2 {
				gameData.publisher = strings.TrimSpace(devAndRating[2])
			}

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

	err := collector.Visit(currentLink)
	if err != nil {
		log.Println("Error when visiting site; pointer set to", err)
		return gameData
	}

	return gameData
}
