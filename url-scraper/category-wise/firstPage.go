package category_wise

import (
	"errors"
	"fmt"
	"github.com/amazon/data-store"
	"github.com/amazon/models"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

const (
	baseUrl      = "https://www.amazon.in"
	storeBaseUrl = baseUrl + "/gp/site-directory?ie=UTF8&ref_=nav_shopall_sbc_fullstore"
)

func ScrapeAllProducts(store data_store.DataStore, category models.Category) error {
	collector := colly.NewCollector()
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Request URL: ", response.Request.URL, " failed with response: ", response, "\nError:", err)
	})

	productCollector := collector.Clone()
	categoryMatched := false
	collector.OnHTML("a.nav_a", func(element *colly.HTMLElement) {
		productCategory := strings.ToLower(strings.TrimSpace(element.Text))
		if productCategory == category.CategoryName {
			fmt.Printf("found category: %v\n", category.CategoryName)
			categoryMatched = true
			catId, err := store.AddCategory(category)
			if err != nil {
				fmt.Printf("error adding category in DB: %v\n", err)
			}
			err = getAllProducts(productCollector, store, catId, baseUrl+element.Attr("href"))
			if err != nil {
				fmt.Printf("error getting list of products: %v\n", err)
			}
		}
	})

	err := collector.Visit(storeBaseUrl)
	if err != nil {
		return err
	}
	if !categoryMatched {
		return errors.New("no such category found")
	}
	return err
}

func getAllProducts(collector *colly.Collector, store data_store.DataStore, catId uint64, productsUrl string) error {
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Request URL:", response.Request.URL, "failed with response:", response, "\nError:", err)
	})

	getDetails(collector, store, catId)

	var numPages string
	collector.OnHTML("span.pagnDisabled", func(element *colly.HTMLElement) {
		numPages = element.Text
	})
	var nextPageLink string
	collector.OnHTML("a#pagnNextLink.pagnNext", func(element *colly.HTMLElement) {
		nextPageLink = element.Attr("href")
		nextPageLink = baseUrl + nextPageLink
	})

	err := collector.Visit(productsUrl)
	if err != nil {
		return err
	}
	totalPages, err := strconv.Atoi(numPages)
	if err != nil {
		fmt.Println("This is the last page.")
		return nil
	}
	err = ScrapeNextPage(totalPages, store, catId, nextPageLink)
	return err
}

func getDetails(collector *colly.Collector, store data_store.DataStore, catId uint64) {
	collector.OnError(func(response *colly.Response, err error) {
		fmt.Println("Request URL:", response.Request.URL, "failed with response:", response, "\nError:", err)
	})

	collector.OnHTML("div#mainResults", func(element *colly.HTMLElement) {
		product := &models.Product{}
		element.ForEach("li[data-asin]", func(i int, element *colly.HTMLElement) {
			product.DataAsin = element.Attr("data-asin")

			element.ForEach("div", func(i int, element *colly.HTMLElement) {
				attr := element.Attr("class")
				if strings.Contains(attr, "a-row a-spacing-none") {
					// product_name:
					element.ForEach("h2", func(i int, element *colly.HTMLElement) {
						product.ProductName = strings.TrimSpace(element.Attr("data-attribute"))
					})

					// product_manufacturer:
					mf := false
					element.ForEachWithBreak("span", func(i int, element *colly.HTMLElement) bool {
						if mf == true {
							product.Manufacturer = strings.TrimSpace(element.Text)
							return false
						}
						if strings.TrimSpace(element.Text) == "by" {
							mf = true
						}
						return true
					})

					// product_price:
					element.ForEach("a", func(i int, element *colly.HTMLElement) {
						element.ForEach("span", func(i int, element *colly.HTMLElement) {
							attr := element.Attr("class")
							if strings.Contains(attr, "price") {
								product.Price = strings.TrimSpace(element.Text)
							}
						})
					})
				}

				// product_ratings:
				element.ForEachWithBreak("i", func(i int, element *colly.HTMLElement) bool {
					element.ForEachWithBreak("span.a-icon-alt", func(i int, element *colly.HTMLElement) bool {
						product.Ratings = strings.TrimSpace(element.Text)
						return false
					})
					return false
				})

				// product_description:
				element.ForEachWithBreak("dl", func(i int, element *colly.HTMLElement) bool {
					desc := ""
					element.ForEachWithBreak("li", func(i int, element *colly.HTMLElement) bool {
						desc += strings.TrimSpace(element.Text) + "\n"
						return true
					})
					product.Description = desc
					return false
				})
			})

			product.CategoryID = catId
			err := store.UpsertProduct(*product)
			if err != nil {
				fmt.Printf("error in saving product in DB: %v\n", err)
			}
		})
	})
}
