package image_fetcher

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Parser interface {
	Parse(url string) (imageSources []string, links []string, err error)
}

type HTMLParser struct{}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

const (
	imgSrcAttr = "img"
	srcKey     = "src"

	aAttr   = "a"
	hrefKey = "href"
)

const (
	httpPrefix  = "http"
	httpsPrefix = "https"
)

// Parse parses the HTML content from the given URL and extracts the image sources (src attributes) of all <img> tags.
func (p *HTMLParser) Parse(urlStr string) (imageSources []string, links []string, err error) {
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	document, err := html.Parse(response.Body)
	if err != nil {
		return nil, nil, err
	}

	base, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	imageSources, links = p.getSources(document, base)

	return imageSources, links, nil
}

// getSources extracts image sources and links from the HTML node.
func (p *HTMLParser) getSources(node *html.Node, base *url.URL) (imgSources []string, linksSources []string) {
	if node.Type == html.ElementNode {
		if node.Data == imgSrcAttr {
			for _, attr := range node.Attr {
				if attr.Key == srcKey {
					src := attr.Val
					if !strings.HasPrefix(src, "http") && !strings.HasPrefix(src, "https") {
						src = resolveURL(src, base)
					}
					imgSources = append(imgSources, src)
				}
			}
		} else if node.Data == aAttr {
			for _, attr := range node.Attr {
				if attr.Key == hrefKey {
					href := attr.Val
					if strings.HasPrefix(href, httpPrefix) || strings.HasPrefix(href, httpsPrefix) {
						linksSources = append(linksSources, href)
					} else {
						href = resolveURL(href, base)
						linksSources = append(linksSources, href)
					}
				}
			}
		}
	}

	for firstChild := node.FirstChild; firstChild != nil; firstChild = firstChild.NextSibling {
		images, links := p.getSources(firstChild, base)

		imgSources = append(imgSources, images...)

		linksSources = append(linksSources, links...)
	}

	return imgSources, linksSources
}

// resolveURL resolves a relative URL to an absolute URL using the base URL.
func resolveURL(ref string, base *url.URL) string {
	u, err := base.Parse(ref)
	if err != nil {
		return ref
	}

	return u.String()
}
