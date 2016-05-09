// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	urlFile := "sites.txt"
	if len(os.Args) > 1 {
		urlFile = os.Args[1]
	}
	urls := make(chan string)
	go urlReader(urlFile, urls)
	pages := make(chan string)
	go pageFetcher(urls, pages)
	done := make(chan bool)
	go pagePrinter(pages, done)
	<-done
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func urlReader(path string, out chan<- string) {
	defer close(out)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func pageFetcher(urls <-chan string, out chan<- string) {
	defer close(out)
	start := time.Now()
	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			out <- fmt.Sprint(err) // send to channel ch
			continue
		}

		file, err := os.Create(fmt.Sprintf("%x", md5.Sum([]byte(url))))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		nbytes, err := io.Copy(file, resp.Body)
		defer resp.Body.Close() // don't leak resources
		if err != nil {
			out <- fmt.Sprintf("while reading %s: %v", url, err)
			continue
		}
		secs := time.Since(start).Seconds()
		out <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
	}
}

func pagePrinter(pages <-chan string, done chan<- bool) {
	for page := range pages {
		fmt.Println(page)
	}
	done <- true
}

//!-
