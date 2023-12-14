package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// fetch {url} and get <title> of the page
// if not found title, return error
func getTitle(url string) (string, error) {
	// fetch url
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// if response status code is not 200, return error
	if resp.StatusCode != 200 {
		return "", errors.New(url + " status code is not 200")
	}
	// extract <title>*</title>
	title, err := getTitleFromResponse(resp)
	if err != nil {
		return "", err
	}
	// trim space
	title = strings.TrimSpace(title)
	return title, nil

}

// extract <h1>*</h1> from response
// if not found, return error
func getTitleFromResponse(resp *http.Response) (string, error) {
	// get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// convert to string
	bodyStr := string(body)
	// pattern match <h1>*</h1>
	titleStart := strings.Index(bodyStr, "<h1>")
	titleEnd := strings.Index(bodyStr, "</h1>")
	if titleStart == -1 || titleEnd == -1 {
		return "", errors.New("Title not found")
	}
	// return title
	return bodyStr[titleStart+4 : titleEnd], nil
}

// main
func main() {
	// comapny id list [1....100000]
	idList := []int{}
	for i := 1; i <= 10_0000; i++ {
		idList = append(idList, i)
	}
	// open companies.csv
	f, err := os.OpenFile("companies.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// if companies.csv is empty, appen header
	// "id,title,url"
	if _, err := f.Stat(); os.IsNotExist(err) {
		if _, err := f.WriteString("id,title,url\n"); err != nil {
			log.Fatal(err)
		}
	}
	// each id fetch
	for _, id := range idList {
		url := fmt.Sprintf("https://k.secure.freee.co.jp/companies/%d/announces", id)
		title, err := getTitle(url)
		if err != nil {
			fmt.Println(err)
		} else {
			// "id: %, comapny: %s, url: %s"
			fmt.Printf("id: %d, comapny: %s, url: %s\n", id, title, url)
			// append "id,title,url" to csv
			// write to companies.csv
			if _, err := f.WriteString(fmt.Sprintf("%d,%s,%s\n", id, title, url)); err != nil {
				log.Fatal(err)
			}
		}
		// sleep 1 sec to avoid too many requests
		time.Sleep(1 * time.Second)
	}
	// clean up
	f.Sync()

}
