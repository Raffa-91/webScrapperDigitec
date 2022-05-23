package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

//  https://www.youtube.com/watch?v=mS74M-rnc90 - Tutorial WebScapper Ebay

func getHtml(url string) *http.Response {
	response, err := http.Get(url)
	checkErr(err)

	if response.StatusCode > 400 {
		fmt.Println("Status Code:", response.StatusCode)
	}

	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}
	return response
}

func scrapeHtml(doc *goquery.Document) {
	doc.Find("div.sc-1f6e68i-0.dkAAXO>article.panelProduct").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.product--title")
		title := a.Text()
		fmt.Println(title) //Why the fuck gits kei HTML Code sondern e pointer uus??
	})
}

func main() {

	url := "https://www.digitec.ch/search?q=maus"

	response := getHtml(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	println(doc)
	checkErr(err)

	scrapeHtml(doc)

}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
