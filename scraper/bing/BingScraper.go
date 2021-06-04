package bingscraper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Dhsieh/tree_scraper/data"
	"github.com/gocolly/colly"
)

func Setup(downloadPath string) *BingScraper {
	newDownloadPath := fmt.Sprintf("%s/bing", downloadPath)
	return &BingScraper{downloadPath: newDownloadPath, counter: 0}
}

func CreateScraper(path string, jsonMap map[string]data.TreeJson) BingScraper {
	return BingScraper{downloadPath: path, treeJsonMap: jsonMap, counter: 0}
}

type BingScraper struct {
	downloadPath string
	treeJsonMap  map[string]data.TreeJson
	counter      int
}

func (b BingScraper) ScrapeImages(treeData data.TreeJson) {
	c := colly.NewCollector()

	tree := treeData.ScientificName
	preparedTree := strings.ReplaceAll(tree, " ", "+")
	url := fmt.Sprintf(bingPhotoUrl, preparedTree)
	fmt.Printf("Scraping %s\n", url)

	var urls []string
	alreadyScrapedCount := 0

	c.OnHTML("li", func(e *colly.HTMLElement) {
		t := e.ChildAttr("a", "m")
		var result map[string]interface{}
		json.Unmarshal([]byte(t), &result)

		murl, ok := result["murl"]
		if ok {
			if !strings.Contains(murl.(string), "bugwood") {
				urls = append(urls, murl.(string))
			} else {
				alreadyScrapedCount++
			}
		}
	})

	c.Visit(url)
	fmt.Printf("Already scraped %d elements\n", alreadyScrapedCount)

	treeName := strings.ReplaceAll(treeData.ScientificName, " ", "_")
	treeName = strings.ReplaceAll(treeName, "/", "_")
	dirName := fmt.Sprintf("%s/%s", b.downloadPath, treeName)
	fmt.Println(dirName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		fmt.Printf("Directy %s was not created, creating it!\n", dirName)
		os.MkdirAll(dirName, os.ModePerm)
	}

	if len(urls) >= 5 {
		b.downloadImages(urls, dirName)
	} else {
		fmt.Printf("Could not find 5 images for %s", tree)
	}
}

func (b BingScraper) downloadImages(urls []string, dir string) {
	counter := 0
	for _, url := range urls {
		b.downloadImage(url, dir, counter)
		counter++
	}
}

func (b BingScraper) downloadImage(url, dir string, counter int) {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		filename := fmt.Sprintf("%s/%d.jpg", dir, counter)

		err := r.Save(filename)
		if err != nil {
			panic(err)
		}
	})

	c.Visit(url)
}

func (b BingScraper) ScrapeAllTrees() {
	for _, tree := range b.treeJsonMap {
		b.ScrapeImages(tree)
	}

}
