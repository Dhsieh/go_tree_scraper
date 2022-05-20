package bingscraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Dhsieh/tree_scraper/data"
	"github.com/Dhsieh/tree_scraper/utils"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Create a BingScraper struct
func Setup(downloadPath string) *BingScraper {
	newDownloadPath := fmt.Sprintf("%s/bing", downloadPath)
	return &BingScraper{downloadPath: newDownloadPath, counter: 0}
}

func CreateScraper(path string, jsonMap map[string]data.TreeJson, images int, numRoutines int, indexFile string, numUrls int) BingScraper {
	return BingScraper{downloadPath: path, treeJsonMap: jsonMap, images: images, counter: 0, numRoutines: numRoutines, indexFile: indexFile, numUrls: numUrls}
}

// Scraper to scrape images from bing
// downloadPath: 	string 						path to save the images
// treeJsonMap: 	map[string]data.TreeJson 	map containing all the species to download
// images: 			int 						number of images to scrape per species
// numRoutines: 	int 						number of go routines to create
// counter: 		int 						way to track number of images saved
// indexFile 	 	string 						path of file to read and update index
// numUrls   		int 						number of search urls to look for images
type BingScraper struct {
	downloadPath string
	treeJsonMap  map[string]data.TreeJson
	images       int
	numRoutines  int
	counter      int
	indexFile    string
	numUrls      int
}

// Determine how to get the images
// There are 2 ways:
//   1. Keyword: Given a keyword scrape images from searching that keyword
//   2. TreeData: List of tree species from bugwood website and the images from each page for each species
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

// Get all the image urls from one bing url
// image urls are from the li tags of the bing url
func (b BingScraper) scrapeImages(url string) []string {
	fmt.Printf("Scraping %s\n", url)
	c := colly.NewCollector()
	noMUrl := 0
	validMUrl := 0
	var urls []string

	c.OnHTML("li", func(e *colly.HTMLElement) {
		tag := e.ChildAttr("a", "m")
		var result map[string]interface{}
		json.Unmarshal([]byte(tag), &result)

		murl, ok := result["murl"]
		if ok {
			urls = append(urls, murl.(string))
			validMUrl++
		} else {
			noMUrl++
		}
	})

	c.Visit(url)
	fmt.Printf("%d urls did not have a murl\n", noMUrl)
	fmt.Printf("valid murl is %d\n", validMUrl)
	return b.checkUrls(urls, b.images)
}

// Scrape images and then store them into GCS
// Create bing urls and download the images from those urls
// bing urls are generated based on an index from a file in GCS
func (b BingScraper) ScrapeImagesToGCS(ctx context.Context, keyword string) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucket := client.Bucket(b.downloadPath)
	index := getIndex(ctx, keyword, bucket)
	urls, index := createUrls(keyword, b.numUrls, index)
	writeIndex(ctx, bucket, index, b.downloadPath, b.indexFile)
	allUrls := b.getAllImageUrls(urls)

	fmt.Printf("The number of urls is %d\n", len(urls))
	b.downloadImages(ctx, allUrls, b.downloadPath)
}

// Function to scrape images of a keyword
func (b BingScraper) ScrapeKeywordImages(keyword string) {
	urls, _ := createUrls(keyword, 3, 0)
	allUrls := b.getAllImageUrls(urls)

	dirName := fmt.Sprintf("%s/%s", b.downloadPath, "trees")
	utils.CreateDirectory(dirName)

	fmt.Printf("Downloading %d images into %s \n", b.images, dirName)

	b.downloadImageList(allUrls, dirName)
}

// Scrape from bugwood
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

// Grab all image urls from the bing urls
func (b BingScraper) getAllImageUrls(urls []string) []string {
	var allUrls []string

	for _, bingUrl := range urls {
		imageUrls := b.scrapeImages(bingUrl)
		allUrls = append(allUrls, imageUrls...)
	}

	return allUrls
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

// Download images using a web service and use channels to parrallelize the images
// urls: 	List of urls
// bucket:  Bucket to download images to
func (b BingScraper) downloadImages(ctx context.Context, urls []string, bucket string) {
	var workerGroup sync.WaitGroup

	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("error creating client! %v", err)
		panic(err)
	}

	in := make(chan string)
	for i := 0; i < b.numRoutines; i++ {
		workerGroup.Add(1)
		go downloadImage(ctx, client, in, bucket, &workerGroup)
	}

	for _, url := range urls {
		in <- url
	}

	close(in)
	workerGroup.Wait()
}

// download an image in a single channel
// TODO: Change name of the files to be less random
// storageClient: GCS storageClient
// in: 			  channel of urls to read from
// dir: 		  bucket in GCS
func downloadImage(ctx context.Context, storageClient *storage.Client, in <-chan string, dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	bucket := storageClient.Bucket(dir)

	for {
		url, ok := <-in
		if !ok {
			break
		}

		c := colly.NewCollector()

		c.OnResponse(func(r *colly.Response) {
			//date format: YYYY_MM_DD
			date := time.Now().Format("2006_01_02")
			id := uuid.New()
			fileName := fmt.Sprintf("download/bing/%s_%s.jpeg", date, id.String())
			file := bucket.Object(fileName)

			writer := file.NewWriter(ctx)
			writer.ContentType = "image/jpeg"
			writer.Name = fileName
			defer writer.Close()

			img, _, err := image.Decode(bytes.NewReader(r.Body))
			if err != nil {
				fmt.Printf("Error decoding byte array to image! %v\n", err)
				panic(err)
			}

			if err := jpeg.Encode(writer, img, nil); err != nil {
				fmt.Printf("Error writing file: %s!", fileName)
				panic(err)
			}

			if err := writer.Close(); err != nil {
				fmt.Printf("Error closing writer %v!", err)
				panic(err)
			}

		})

		c.Visit(url)
	}
}

// Create a list of urls where each url will contain 35 images
// keyword: Keyword to place into URL
// numUrls: number of urls to scrape on each page
// index: number of images on the url, the next page will be index*35 + 1
func createUrls(keyword string, numUrls int, index int) ([]string, int) {
	var urls []string
	for i := 0; i < numUrls; i, index = i+1, index+1 {
		newUrl := fmt.Sprintf(bingKeywordPhotoNumberUrl, keyword, index*35+1)
		urls = append(urls, newUrl)
	}
	return urls, index
}

// Read index file to see where to start when creating Bing urls
func getIndex(ctx context.Context, keyword string, bucket *storage.BucketHandle) int {

	attrs, _ := bucket.Attrs(ctx)
	dir := attrs.Name
	var indexFile *storage.ObjectHandle
	index := 0

	it := bucket.Objects(ctx, nil)
	// Find the indexfile to see if it exists
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			fmt.Println("Could not find the indexfile!")
			break
		}
		if err != nil {
			fmt.Printf("Bucket(%q).Objects: %v", dir, err)
		}
		if strings.Contains(attrs.Name, "indexfile.txt") {
			fmt.Printf("File name is %s\n", attrs.Name)
			indexFile = bucket.Object("index/indexfile.txt")
			//break once the correct file is found.
			break
		}
	}

	if indexFile == nil {
		fmt.Println("There was no index file!")
	} else {
		// Read the contents of the file and return
		fmt.Println("Reading file contents now!")
		reader, err := indexFile.NewReader(ctx)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer reader.Close()

		fileContents, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		index, err = strconv.Atoi(string(fileContents))
		fmt.Printf("Index is %d\n", index)
	}

	// default value is 0
	return index
}

// Write index file to a file in GCS
func writeIndex(ctx context.Context, bucket *storage.BucketHandle, index int, dir string, file string) {

	fmt.Printf("Writing to file %s\n", file)
	indexFile := bucket.Object(file)
	writer := indexFile.NewWriter(ctx)
	writer.ContentType = "text/plain"

	if _, err := writer.Write([]byte(strconv.Itoa(index))); err != nil {
		fmt.Printf("createFile: unable to write data to bucket %q, file %q: %v\n", dir, file, err)
		panic(err)
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("Could not close file!\n")
		panic(err)
	}
}

// Print BingScraper struct
func (b BingScraper) String() string {
	return fmt.Sprintf("downloadPath: %s\ntreeJsonMap: %v\nNumber of Images: %d\nNumber of Routines %d\nCounter %d", b.downloadPath, b.treeJsonMap, b.images, b.numRoutines, b.counter)
}
