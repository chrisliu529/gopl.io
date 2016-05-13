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
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	urlFile := "sites.txt"
	if len(os.Args) > 1 {
		urlFile = os.Args[1]
	}
	processUrls(urlFile, handlePage)
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func processUrls(path string,
	handler func(string, chan<- string, *sync.WaitGroup)) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	pages := make(chan string, 256)
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go handler(scanner.Text(), pages, &wg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	go func() {
		wg.Wait()
		close(pages)
	}()

	for page := range pages {
		fmt.Println(page)
	}
}

func handlePage(url string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		out <- fmt.Sprint(err)
		return
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
		return
	}
	secs := time.Since(start).Seconds()
	out <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//!-
