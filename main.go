package main

import (
	"github.com/amazon/data-store"
	"github.com/amazon/models"
	"github.com/amazon/url-scraper/category-wise"
	"github.com/amazon/url-scraper/single-product"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

const defaultPort = "8000"

var store data_store.DataStore

func productInfo(c *gin.Context) {
	link := c.Request.FormValue("url")
	err := single_product.ScrapeProductInfo(store, link)
	if err != nil {
		log.Printf("[Scraper Error]: %v\n", err)
	}
	_, err = c.Writer.WriteString("Product Scraping: ok")
	if err != nil {
		log.Println(err)
	}
}

func getAllProducts(c *gin.Context) {
	category := &models.Category{
		CategoryName: c.Request.FormValue("category"),
	}
	category.CategoryName = strings.ToLower(category.CategoryName)
	err := category_wise.ScrapeAllProducts(store, *category)
	if err != nil {
		log.Printf("[Scraper Error]: %v\n", err)
	}
	_, err = c.Writer.WriteString("Category-wise Scraping: ok")
	if err != nil {
		log.Println(err)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	store = data_store.DbConnect()
	router := gin.Default()
	router.POST("/product", productInfo)
	router.POST("/all-products", getAllProducts)
	log.Println("server up...")
	log.Fatal(http.ListenAndServe(":"+defaultPort, router))
}
