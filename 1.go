package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/html"
)

func main() {
	var (
		k       = 5
		counter int64
	)

	limitChan := make(chan struct{}, k)
	wg := &sync.WaitGroup{}

	var url string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		limitChan <- struct{}{}
		url = scanner.Text()
		go analyze(limitChan, wg, url, &counter)
	}
	wg.Wait()
	fmt.Printf("Total: %d\n", counter)
}

// analyze общая функция для получения данных с сайта и парсинга контента
func analyze(limitChan chan struct{}, wg *sync.WaitGroup, url string, counter *int64) {
	content, err := getContent(url)
	if err != nil {
		log.Println("Something wrong with", url)
	} else {
		parsedContent := parseContent(content)
		count := strings.Count(parsedContent, `Go`)
		atomic.AddInt64(counter, int64(count))
		fmt.Printf("Count for %s: %d\n", url, count)
	}
	<-limitChan
	wg.Done()
}

// getContent возвращает полученный контент с указанного url
func getContent(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(html), nil
}

// parseContent возвращает текст из html контента
func parseContent(content string) (parsedContent string) {
	domDocTest := html.NewTokenizer(strings.NewReader(content))
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" || previousStartTokenTest.Data == "title" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			if len(TxtContent) > 0 {
				parsedContent += TxtContent
			}
		}
	}
	return
}
