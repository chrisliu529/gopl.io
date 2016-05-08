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
	urls := readLines(urlFile)
	pages := make(chan string)
	for _, url := range urls {
		go fetch(url, pages) // start a goroutine
	}
	for range urls {
		fmt.Println(<-pages)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func readLines(path string) []string {
	lines := []string{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
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
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	res := fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
	ch <- res
}

//!-
