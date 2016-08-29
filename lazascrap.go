package main

import "lazada-scrapper/lazascrap"
import (
	"fmt"
)

const Version = "0.1.0"
const BaseUrl = "http://www.lazada.co.id/promo-hsbc/"

func main() {
	product, err := lazascrap.ScrapPage(BaseUrl, 1)

	if err == nil {
		fmt.Println("Total: ", product.TotalItems)
		for i, item := range product.Items {
			fmt.Println("Product", (i + 1))
			fmt.Println("\tTitle:", item.Title)
			fmt.Println("\tOriginal Price:", item.Price)
			fmt.Println("\tDiscounted Price:", item.DiscountedPrice)
			fmt.Println("\tImage:", item.Image)
		}
	} else {
		fmt.Println(err)
	}

}
