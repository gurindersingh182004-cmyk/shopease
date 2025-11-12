package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	pageURL := "https://example.com" // 🔹 Change this to the website you want to scrape

	c := colly.NewCollector()

	c.OnHTML("img", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		if imgSrc == "" {
			return
		}

		// Resolve relative URLs
		u, err := url.Parse(imgSrc)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		base, _ := url.Parse(pageURL)
		imgURL := base.ResolveReference(u).String()

		fmt.Println("Downloading:", imgURL)
		downloadImage(imgURL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	fmt.Println("Visiting:", pageURL)
	c.Visit(pageURL)
}

// downloadImage downloads an image from the URL and saves it locally
func downloadImage(imgURL string) {
	resp, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("Download error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Non-200 status:", resp.Status)
		return
	}

	filename := path.Base(resp.Request.URL.Path)
	if filename == "" || strings.Contains(filename, "?") {
		filename = "image.jpg"
	}

	out, err := os.Create(filename)
	if err != nil {
		fmt.Println("File create error:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error saving:", err)
		return
	}

	fmt.Println("Saved:", filename)
}
