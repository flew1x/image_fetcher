package image_fetcher

import (
	"github.com/flew1x/image_fetcher/utils"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Downloader interface {
	Download(src string, baseURL string, path string) error
}

type HTTPDownloader struct{}

func NewHTTPDownloader() *HTTPDownloader {
	return &HTTPDownloader{}
}

// Download downloads an image from the specified source URL.
//
// Parameters:
// - src: the URL of the image to download.
// - baseURL: the base URL for resolving relative URLs.
// - path: the directory path to save the downloaded image.
// Returns:
// - error: an error if there was an issue during the download process.
func (d *HTTPDownloader) Download(source string, baseURL string, path string) error {
	imageURL, err := d.getResolvedReference(source, baseURL)
	if err != nil {
		return err
	}

	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	fileName := filepath.Base(imageURL)
	filePath := filepath.Join(path, fileName)

	return utils.SaveToFile(response.Body, filePath)
}

// getResolvedReference resolves the image URL based on the base URL.
//
// Parameters:
// - imageURL: the URL of the image to resolve.
// - baseURL: the base URL used for resolving relative URLs.
// Returns:
// - string: the resolved image URL.
// - error: an error if there was an issue during the resolution process.
func (d *HTTPDownloader) getResolvedReference(imageURL string, baseURL string) (string, error) {
	if !strings.HasPrefix(imageURL, httpPrefix) {
		base, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}

		ref, err := url.Parse(imageURL)
		if err != nil {
			return "", err
		}

		imageURL = base.ResolveReference(ref).String()
	}

	return imageURL, nil
}
