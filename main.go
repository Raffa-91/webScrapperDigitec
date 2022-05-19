package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	url := "https://www.digitec.ch/search?q=maus"
	code := getCode(url)

	file, err := os.Create("htmlCode.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(code)
	println("HTML Code was written to file")

	/*
		stringStart := strings.Index(code, "class=\"sc-qlvix8-0 dCvedB\" aria-label=")
		if stringStart == -1 {
			fmt.Println("No Element Found")
			os.Exit()
		}
		stringStart += 2

		stringEnd := strings.Index(code, "/"
		"")
		if stringEnd == -1 {
			fmt.Println("No Element Found")
			os.Exit()
		}
	*/
}

func getCode(url string) string {

	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	result := html
	return string(result)
}
