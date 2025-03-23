package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"golang.org/x/net/html"
)

func main() {
	var urls []string
    fmt.Println("Enter URLs (type 'exit' to stop):")
	for {
		var url string
		fmt.Scan(&url)
		if url == "exit" {
			break
		}
		urls = append(urls, url)
	}
    log.Println("Extracting Links form:",urls)

	linksChannel := make(chan []string)
	var wg sync.WaitGroup

	for _,url := range urls {
		wg.Add(1)
		go extractLinks(url, linksChannel, &wg)
	}

	for links := range linksChannel {
		fmt.Println("Links found:", links)
	}

	wg.Wait()
}

func extractLinks(url string, linksChannel chan []string,wg *sync.WaitGroup)  {
	res, _ := http.Get(url)
	var links []string

	z := html.NewTokenizer(res.Body)
	defer res.Body.Close()

	for {
		tt := z.Next()


		if tt == html.ErrorToken {
			break
		}
		t:= z.Token()
		isAnchor := t.Data == "a"
 
		if isAnchor {
			for _,a := range t.Attr {
				if a.Key == "href" && !strings.Contains(a.Val, url)  && strings.Contains(a.Val, "https") {
					links = append(links, a.Val)
				}
			}
		}

	}
	wg.Done()
	linksChannel <- links
}
