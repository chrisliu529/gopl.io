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
	processUrls(urlFile, handlePage)
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func processUrls(path string, handler func(string, chan<- bool)) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	scanner := bufio.NewScanner(file)
	n := 0
	for scanner.Scan() {
		go handler(scanner.Text(), done)
		n++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < n; i++ {
		<-done
	}
}

func handlePage(url string, done chan<- bool) {
	h := func(s string) {
		fmt.Println(s)
		done <- true
	}
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		h(fmt.Sprint(err))
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
		h(fmt.Sprintf("while reading %s: %v", url, err))
		return
	}
	secs := time.Since(start).Seconds()
	h(fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url))
}

//!-
