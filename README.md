#amazon-scraper

amazon-scraper is build for performing two types of scraping on Amazon platform. Given product url, it can scrape the details of the product. It also supports the category wise listing of all the products on Amazon and stores the all product details in the database. 

## To run the scraper:

####Step-1:
 Clone the repository and acquire all dependencies. 
 
 `go mod download`

####Step-2:
 Setup MySQL and create database with database name as `amazon`.
 
####Step-3:
 Run the server using the command:
 
`go run main.go`

There are two endpoints for scraping the information from Amazon:

####`/product` API:
It takes the url of the product as input and provides the information about the product such as:
`productName`
`category_id`
`manufacturer`
`price`
`ratings`
`description`

#### `/all-products` API:

It takes `category-name` as input and provides all the products in the specified category. It stores all the products with the detailed information in the database. It will scrape the products continuously till the last item in the category.  