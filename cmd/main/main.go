package main

import (
	"flag"
	"fmt"
	"github.com/flew1x/image_fetcher"
	"log"
	"os"
	"time"
)

const (
	// URL to crawl
	urlText = "URL to crawl"

	// Path to save images
	pathText = "Path to save images"

	// Depth for internal links
	depthText = "Depth for internal links"
)

func main() {
	// Parse flags
	urlFlag := flag.String("f", "", urlText)
	pathFlag := flag.String("p", "", pathText)
	depthFlag := flag.Int("e", 0, depthText)
	timeoutFlag := flag.String("t", "10m", "Timeout in seconds")
	flag.Parse()

	// Check if flags are set
	if *urlFlag == "" || *pathFlag == "" || *depthFlag < 0 || *timeoutFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	parsedDelay, err := time.ParseDuration(*timeoutFlag)
	if err != nil {
		log.Fatalf("Error parsing timeout: %v\n", err)
	}

	// Create directory
	err = os.MkdirAll(*pathFlag, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
	}

	parser := image_fetcher.NewHTMLParser()
	downloader := image_fetcher.NewHTTPDownloader()

	crawler := image_fetcher.NewCrawler(parser, downloader, *pathFlag, parsedDelay)

	err = crawler.Crawl(*urlFlag, *depthFlag)
	if err != nil {
		log.Fatalf("Error during crawling: %v\n", err)
	}

	fmt.Println("Crawling and downloading completed.")
}
