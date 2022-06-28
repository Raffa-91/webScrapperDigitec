package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	_ "go/printer"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

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

func saveToCSV(result []string) {
	file, error := os.Create("results.csv")
	checkErr(error)
	writer := csv.NewWriter(file)
	writer.Write(result)
	writer.Flush()
}

func scrapeHtml(doc *goquery.Document) []string {
	var result []string

	doc.Find(".product-wrapper").Each(func(index int, item *goquery.Selection) {
		a := item.Find("h2 a")
		price := item.Find(".price.small")
		title := strings.TrimSpace(a.Text())
		priceValue := strings.TrimSpace(price.Text())
		fmt.Println("Titel:")
		fmt.Printf("%+v \n", title)
		fmt.Println("Preis:")
		fmt.Printf("%+v \n", priceValue)
		result = append(result, title, priceValue, "\n")
	})

	return result
}

func main() {
	ClearTerminal()
	PrintMenue()
	for true {
		ExecuteCommand()
	}

}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func searchThroughWebsite(page int, command string, result *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://www.mediamarkt.ch/de/search.html?query=%s&searchProfile=onlineshop&channel=mmchde&page=%d", command, page)

	response := getHtml(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	fmt.Println(&doc)
	checkErr(err)

	*result = append(*result, scrapeHtml(doc)...)
}

func loopOver(command string) []string {
	var result []string

	var wg sync.WaitGroup
	wg.Add(5)
	//Markenbezeichnungen in "command" generienen teilweise Unterseiten welche nicht gelesen werden können
	go searchThroughWebsite(1, command, &result, &wg)
	go searchThroughWebsite(2, command, &result, &wg)
	go searchThroughWebsite(3, command, &result, &wg)
	go searchThroughWebsite(4, command, &result, &wg)
	go searchThroughWebsite(5, command, &result, &wg)
	wg.Wait()

	// Save results into a CSV File
	saveToCSV(result)

	return result
}

func PrintMenue() {
	fmt.Println(`
##########################################
#******** Webscapper Mediamarkt **********
#********* Wähle eine Option *************
# 1. Suchen
# 2. Export (Beta)
#
# q. Beenden
##########################################
##########################################
##########################################
`)
}

func askForCommand() string {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.TrimSpace(response)
	return response
}

func ExecuteCommand() {
	command := askForCommand()
	parseCommand(command)
}

func ClearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func parseCommand(input string) {
	switch {
	case input == "1":
		ClearTerminal()
		fmt.Println("Suchwort eingeben:")
		command := askForCommand()
		loopOver(command)
		break
	case input == "2":
		ClearTerminal()
		// Sample für (Beta)Export
		loopOver("Tastatur")
		PrintMenue()
		break
	}
}
