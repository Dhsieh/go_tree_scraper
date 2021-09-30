package bingscraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dhsieh/tree_scraper/data"
	"github.com/Dhsieh/tree_scraper/utils"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

// Create a BingScraper struct
func Setup(downloadPath string) *BingScraper {
	newDownloadPath := fmt.Sprintf("%s/bing", downloadPath)
	return &BingScraper{downloadPath: newDownloadPath, counter: 0}
}

func CreateScraper(path string, jsonMap map[string]data.TreeJson, images int, numRoutines int) BingScraper {
	return BingScraper{downloadPath: path, treeJsonMap: jsonMap, images: images, counter: 0, numRoutines: numRoutines}
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

func (b BingScraper) ScrapeImages(input interface{}) {
	switch in := input.(type) {
	case string:
		fmt.Println("Scraping keyword!")
		b.ScrapeKeywordImages(in)
	case []data.TreeJson:
		b.ScrapeTreeDatas(in)
	default:
		panic("Could not create a function for give type!")
	}
}

func (b BingScraper) scrapeImages(url string) []string {
	fmt.Printf("Scraping %s", url)
	c := colly.NewCollector()
	var urls []string

	c.OnHTML("li", func(e *colly.HTMLElement) {
		tag := e.ChildAttr("a", "m")
		var result map[string]interface{}
		json.Unmarshal([]byte(tag), &result)

		murl, ok := result["murl"]
		if ok {
			urls = append(urls, murl.(string))
		}
	})

	c.Visit(url)
	return b.checkUrls(urls, b.images)
}

// Function to scrape images of a keyword
func (b BingScraper) ScrapeKeywordImages(keyword string) {
	url := fmt.Sprintf(bingKeywordPhotoUrl, keyword)
	urls := b.scrapeImages(url)

	dirName := fmt.Sprintf("%s/%s", b.downloadPath, "trees")
	utils.CreateDirectory(dirName)

	fmt.Printf("Downloading %d images into %s \n", b.images, dirName)

	b.downloadImageList(urls, dirName)
}

func (b BingScraper) ScrapeTreeDatas(treeSlice []data.TreeJson) {
	var waitGroup sync.WaitGroup
	in := make(chan data.TreeJson, b.numRoutines*4)

	for i := 0; i < b.numRoutines; i++ {
		waitGroup.Add(1)
		go b.ScrapeTreeData(in, &waitGroup)
	}

	for _, treeData := range treeSlice {
		in <- treeData
	}

	close(in)
	waitGroup.Wait()
}

// Function that will find images for a tree and downloads them
func (b BingScraper) ScrapeTreeData(in chan data.TreeJson, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector()

	for {
		treeData, ok := <-in
		if !ok {
			break
		}
		fmt.Printf("Scraping %s\n", treeData.ScientificName)
		preparedTree := strings.ReplaceAll(treeData.ScientificName, " ", "+")
		url := fmt.Sprintf(bingPhotoUrl, preparedTree)

		var urls []string
		c.OnHTML("li", func(e *colly.HTMLElement) {
			t := e.ChildAttr("a", "m")
			var result map[string]interface{}
			json.Unmarshal([]byte(t), &result)

			murl, ok := result["murl"]
			if ok {
				urls = append(urls, murl.(string))
			} else {
				fmt.Println("Could not find the murl!")
			}
		})

		c.Visit(url)
		if treeData.ScientificName == "" {
			fmt.Printf("ScientificName not found for %s\n", treeData.CommonName)
			return
		}
		treeName := strings.ReplaceAll(treeData.ScientificName, " ", "_")
		treeName = strings.ReplaceAll(treeName, "/", "_")
		treeName = strings.ReplaceAll(treeName, "-", "_")
		dirName := fmt.Sprintf("%s/%s", b.downloadPath, treeName)

		utils.CreateDirectory(dirName)
		urls = b.checkUrls(urls, b.images)
		b.downloadImageList(urls, dirName)
	}

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
// use random string generator to give each image a unique name.
func (b BingScraper) downloadImageList(urls []string, dir string) {
	var workerGroup sync.WaitGroup

	in := make(chan string, b.numRoutines*4)

	fmt.Printf("Num of rounites is %d\n", b.numRoutines)
	for i := 0; i < b.numRoutines; i++ {
		workerGroup.Add(1)
		go download(in, dir, &workerGroup)
	}

	fmt.Printf("Going through %d urls\n", len(urls))
	for _, url := range urls {
		in <- url
	}

	close(in)
	workerGroup.Wait()
}

// Download an image from a url, and the name of the resulting image will be based on the counter
func download(in <-chan string, dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		url, ok := <-in
		if !ok {
			break
		}
		c := colly.NewCollector()

		c.OnResponse(func(r *colly.Response) {
			date := time.Now().Format("2006_01_02")
			id := uuid.New()
			filePath := fmt.Sprintf("%s/%s_%s.jpg", dir, date, id.String())

			err := r.Save(filePath)
			if err != nil {
				panic(err)
			}
		})

		c.Visit(url)
	}
}
