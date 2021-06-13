package bingscraper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Dhsieh/tree_scraper/data"
	"github.com/Dhsieh/tree_scraper/utils"
	"github.com/gocolly/colly"
)

// Create a BingScraper struct
func Setup(downloadPath string) *BingScraper {
	newDownloadPath := fmt.Sprintf("%s/bing", downloadPath)
	return &BingScraper{downloadPath: newDownloadPath, counter: 0}
}

func CreateScraper(path string, jsonMap map[string]data.TreeJson, images int) BingScraper {
	return BingScraper{downloadPath: path, treeJsonMap: jsonMap, images: images, counter: 0}
}

// Scraper to scrape images from bing
// downloadPath: 	string 						path to save the images
// treeJsonMap: 	map[string]data.TreeJson 	map containing all the species to download
// images: 			int 						number of images to scrape per species
// numRoutines: 	int 						number of go routines to create
// counter: 		int 						way to track number of images saved
type BingScraper struct {
	downloadPath string
	treeJsonMap  map[string]data.TreeJson
	images       int
	numRoutines  int
	counter      int
}

// Function that will find images for a tree and download them
func (b BingScraper) ScrapeImages(treeData data.TreeJson) {
	c := colly.NewCollector()

	//Prepare the bing url which uses + instead of spaces
	// Using the scientific name to make sure that the images are the most relevant
	tree := treeData.ScientificName
	preparedTree := strings.ReplaceAll(tree, " ", "+")
	url := fmt.Sprintf(bingPhotoUrl, preparedTree)
	fmt.Printf("Scraping %s\n", url)

	var urls []string
	alreadyScrapedCount := 0

	// Don't scrape any images that could be found in the forestryscraper as to not duplicated images
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

	// Create directory to save the images to
	if treeData.ScientificName == "" {
		fmt.Printf("ScientificName not found for %s\n", treeData.CommonName)
		return
	}

	treeName := strings.ReplaceAll(treeData.ScientificName, " ", "_")
	treeName = strings.ReplaceAll(treeName, "/", "_")
	dirName := fmt.Sprintf("%s/%s", b.downloadPath, treeName)
	fmt.Println(dirName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		fmt.Printf("Directory %s was not created, creating it!\n", dirName)
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Printf("Could not create directory %s\n", dirName)
		}
	} else {
		fmt.Printf("Directory %s was already created!\n", dirName)
	}

	// Only scrape the specified number of images
	fmt.Printf("Downloading %d images into %s \n", b.images, dirName)

	b.downloadImages(b.checkUrls(urls, b.images), dirName)
}

// Check if a url contains a JPEG file or not
func (b BingScraper) isJPEG(url string) bool {
	c := colly.NewCollector()

	// bool to see if the url contains a valid jpeg image or not
	var check bool
	c.OnResponse(func(r *colly.Response) {
		check = utils.IsJPEG(r.Body)
	})

	c.Visit(url)
	if !check {
		fmt.Printf("%s did not contain a jpg image!\n", url)
	}

	return check
}

// For the given urls, grab numImages urls that have valid jpeg
func (b BingScraper) checkUrls(urls []string, numImages int) []string {
	var validJPEGUrls []string
	counter := 0
	for i := 0; counter < numImages && i < len(urls); i += 1 {
		url := urls[i]
		if b.isJPEG(url) {
			validJPEGUrls = append(validJPEGUrls, url)
			counter += 1
		}
	}
	return validJPEGUrls
}

// This could also be changed to use channels and go routines
// Given a list of bing urls, download images from it.
// use counter to give each image a unique name.
func (b BingScraper) downloadImages(urls []string, dir string) {
	counter := b.counter
	for _, url := range urls {
		b.downloadImage(url, dir, counter)
		counter++
	}
}

// Download an image from a url, and the name of the resulting image will be based on the counter
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

func (b BingScraper) ScrapeTree(in <-chan data.TreeJson, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		tree, ok := <-in
		if !ok {
			break
		}
		b.ScrapeImages(tree)
	}
}
