package category_wise

import (
	"fmt"
	"github.com/amazon/data-store"
	"github.com/amazon/models"
	"github.com/gocolly/colly"
	"strings"
)

func ScrapeNextPage(totalPages int, store data_store.DataStore, catId uint64, nextPageLink string) error {
	for i := 0; i < totalPages-1; i++ {
		collector := colly.NewCollector()
		collector.OnRequest(func(request *colly.Request) {
			fmt.Println("Visiting", request.URL)
		})
		collector.OnError(func(response *colly.Response, err error) {
			fmt.Println("Request URL:", response.Request.URL, "failed with response:", response, "\nError:", err)
		})
		collector.OnHTML("div[data-asin]", func(element *colly.HTMLElement) {
			product := &models.Product{}
			product.DataAsin = element.Attr("data-asin")

			// product_name:
			element.ForEach("h2", func(i int, element *colly.HTMLElement) {
				element.ForEach("a", func(i int, element *colly.HTMLElement) {
					element.ForEach("span", func(i int, element *colly.HTMLElement) {
						product.ProductName = element.Text
					})
				})
			})

			// product_price:
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				element.ForEachWithBreak("span.a-price-whole", func(i int, element *colly.HTMLElement) bool {
					product.Price = strings.TrimSpace(element.Text)
					return false
				})
			})

			// product_ratings:
			element.ForEachWithBreak("i", func(i int, element *colly.HTMLElement) bool {
				element.ForEachWithBreak("span.a-icon-alt", func(i int, element *colly.HTMLElement) bool {
					product.Ratings = element.Text
					return false
				})
				return false
			})
			product.CategoryID = catId
			err := store.UpsertProduct(*product)
			if err != nil {
				fmt.Printf("error in saving product in DB: %v\n", err)
			}
		})
		var nextPage string
		collector.OnHTML("li.a-last", func(element *colly.HTMLElement) {
			nextPage = element.ChildAttr("a", "href")
			nextPage = baseUrl + nextPage
		})
		err := collector.Visit(nextPageLink)
		if err != nil {
			return err
		}
		nextPageLink = nextPage
	}
	return nil
}
