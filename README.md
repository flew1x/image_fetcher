# Web Crawler

A simple web crawler written in Go that downloads images from web pages and supports configurable request timeouts and depth limits.

## Features

- **Crawl Web Pages:** Recursively crawls web pages starting from a given URL.
- **Download Images:** Downloads images found on the web pages.
- **Configurable Request Delay:** Allows setting a delay between requests to avoid spamming.
- **Depth Limit:** Controls the depth of the crawl to avoid excessively deep traversals.

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/flew1x/image_fetcher
   cd image_fetcher

# Usage

Run the application from the command line with the required flags. The application accepts the following flags:

`-f` The starting URL for the crawl. <br>
`-p`: The directory path where downloaded images will be saved. <br>
`-e` or --depth: The maximum depth of the crawl. <br>
`-t` or --timeout: The delay between requests, specified as a duration (e.g., 10s, 1m).

Example
To run the web crawler with a 2-minute delay between requests, crawling from http://example.com, saving images to ./images, and with a maximum depth of 3:

```sh
go run main.go -f http://example.com -p ./images -e 3 -t 2m
```

Also you can build the app with `make build`