package image_fetcher

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"
)

type Crawler struct {
	parser     Parser
	downloader Downloader
	path       string
	delay      time.Duration
}

func NewCrawler(p Parser, d Downloader, path string, delay time.Duration) *Crawler {
	return &Crawler{
		parser:     p,
		downloader: d,
		path:       path,
		delay:      delay,
	}
}

// Crawl performs a breadth-first crawl of the web starting from the given URL.
func (c *Crawler) Crawl(startURL string, maxDepth int) error {
	if _, err := url.ParseRequestURI(startURL); err != nil {
		return fmt.Errorf("invalid start URL: %w", err)
	}

	if maxDepth < 0 {
		return fmt.Errorf("maxDepth must be non-negative")
	}

	var wg sync.WaitGroup
	visited := &sync.Map{}

	wg.Add(1)

	go c.crawlPage(startURL, 0, maxDepth, visited, &wg)

	wg.Wait()

	return nil
}

// crawlPage performs a breadth-first crawl of the web starting from the given URL.
func (c *Crawler) crawlPage(url string, depth int, maxDepth int, visited *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth >= maxDepth {
		log.Printf("Reached max depth (%d) for URL: %s\n", maxDepth, url)

		return
	}

	if _, loaded := visited.LoadOrStore(url, true); loaded {
		log.Printf("URL already visited: %s\n", url)

		return
	}

	log.Printf("Processing URL: %s at depth %d\n", url, depth)

	imageSources, links, err := c.parser.Parse(url)
	if err != nil {
		log.Printf("Error parsing URL %s: %v\n", url, err)

		return
	}

	if err := c.downloadImages(imageSources, url); err != nil {
		log.Printf("Error downloading images for URL %s: %v\n", url, err)

		return
	}

	log.Printf("Processed URL: %s Depth: %d Images: %d Links: %d\n", url, depth, len(imageSources), len(links))

	for _, link := range links {
		if link == "" {
			continue
		}

		if _, loaded := visited.Load(link); !loaded {
			wg.Add(1)
			go c.crawlPage(link, depth+1, maxDepth, visited, wg)
		}
	}
}

// downloadImages downloads the images from the given URLs.
func (c *Crawler) downloadImages(imageSources []string, baseURL string) error {
	if len(imageSources) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(imageSources))

	for _, src := range imageSources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()

			c.waitTimeout()

			if err := c.downloader.Download(src, baseURL, c.path); err != nil {
				errs <- fmt.Errorf("error downloading file %s: %w", src, err)
			}
		}(src)
	}

	wg.Wait()
	close(errs)

	var combinedError error
	for err := range errs {
		if combinedError == nil {
			combinedError = err
		} else {
			combinedError = fmt.Errorf("%w; %v", combinedError, err)
		}
	}

	return combinedError
}

// waitTimeout waits for the crawler's timeout period.
func (c *Crawler) waitTimeout() {
	time.Sleep(c.delay)
}
