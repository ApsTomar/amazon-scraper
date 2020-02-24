package single_product

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	data_store "github.com/amazon/data-store"
	"github.com/amazon/models"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
	"strings"
)

func ScrapeProductInfo(store data_store.DataStore, link string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Status code: %v", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error(err)
	}
	product, err := getDetails(doc)
	if err != nil {
		return err
	}
	fmt.Printf("Product info:\n%+v\n", product)
	//err = store.UpsertProduct(*product)
	return err

}

func getDetails(doc *goquery.Document) (*models.ViewProduct, error) {
	name := doc.Find("span#productTitle").Text()
	if name == "" {
		return nil, fmt.Errorf("product_name not found")
	}
	name = strings.TrimSpace(name)

	price := doc.Find("span#priceblock_ourprice").Text()
	if price == "" {
		price = doc.Find("span#priceblock_saleprice").Text()
		if price == "" {
			price = doc.Find("span#priceblock_dealprice").Text()
		}
	}
	if price == "" {
		return nil, fmt.Errorf("product_price not found")
	}
	price = strings.TrimSpace(price)

	desc := ""
	doc.Find("div#feature-bullets").Each(func(i int, selection *goquery.Selection) {
		selection.Find("li").Each(func(i int, selection *goquery.Selection) {
			text := selection.Find("span.a-list-item").Text()
			text = strings.TrimSpace(text)
			desc += text + "\n"
		})
	})
	if desc == "" {
		return nil, fmt.Errorf("no product_description found")
	}

	manufacturer := doc.Find("a#bylineInfo").Text()
	if manufacturer == "" {
		fmt.Println("manufacturer not found")
	}
	manufacturer = strings.TrimSpace(manufacturer)

	var ratings, dataAsin string
	var ok bool
	doc.Find("div#averageCustomerReviews").Each(func(i int, selection *goquery.Selection) {
		dataAsin, ok = selection.Attr("data-asin")
		if !ok {
			fmt.Println("data-asin not found")
		}
		ratings = selection.Find("span.a-icon-alt").Text()
	})
	if ratings == "" {
		return nil, fmt.Errorf("ratings not found")
	}
	product := &models.ViewProduct{
		DataAsin:     dataAsin,
		ProductName:  name,
		Price:        price,
		Description:  desc,
		Manufacturer: manufacturer,
		Ratings:      ratings,
	}
	return product, nil
}
