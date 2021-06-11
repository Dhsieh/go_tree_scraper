package forestryscraper

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Dhsieh/tree_scraper/data"
	"github.com/gocolly/colly"
)

// Save to json file
func downloadJson(treeData *data.TreeResponse) {
	data, _ := json.MarshalIndent(treeData, "", " ")
	filePath := "../downloads/tree_data.json"

	var flag int = os.O_WRONLY | os.O_CREATE

	f, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(treeData.Data))
	if _, err := f.Write(data); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

// Downloads all trees from forestryscraper
func DownloadAllTreeSpecies() {
	fmt.Println("Downloading Conifers")
	treeResponses := getSpecies(coniferTreeListUrl, "confier")
	fmt.Printf("Len of conifers is %d\n", len(treeResponses.Data))
	fmt.Println("Downloading hardwood")
	hardwoodSpecies := getSpecies(hardwoodtreeListUrl, "deciduous")
	fmt.Printf("Len of deciduous is %d\n", len(hardwoodSpecies.Data))
	treeResponses.Append(&hardwoodSpecies)

	downloadJson(&treeResponses)
}

// Gets all the trees from a url
func getSpecies(url, treeType string) data.TreeResponse {
	c := colly.NewCollector()

	treeResponseStruct := data.TreeResponse{}
	c.OnResponse(func(r *colly.Response) {
		treeResponseStruct.UnmarshalJSON(r.Body, treeType)
	})

	c.Visit(url)

	return treeResponseStruct
}

func CreateScraper(path string, jsonMap map[string]data.TreeJson) ForestryImageScraper {
	return ForestryImageScraper{downloadPath: path, treeJsonMap: jsonMap}
}

type ForestryImageScraper struct {
	treeJsonMap  map[string]data.TreeJson
	downloadPath string
}

// Scrapes images from forestryimage
func (f ForestryImageScraper) ScrapeImages(treeData data.TreeJson) {
	tree := treeData.CommonName
	treeJson, ok := f.treeJsonMap[tree]
	if !ok {
		fmt.Errorf("Could not find the common name: %s in the map!", tree)
	}
	catId := treeJson.CategoryId

	toScrape := fmt.Sprintf(plantImageListUrl, catId)
	c := colly.NewCollector(
		colly.UserAgent("Test of BugWoodAPI!"),
	)

	// Get all the images
	imageResponseStruct := data.ImageResponse{}
	c.OnResponse(func(r *colly.Response) {
		json.Unmarshal(r.Body, &imageResponseStruct)
	})

	c.Visit(toScrape)
	fmt.Printf("%s has %d images to scrape\n", tree, imageResponseStruct.Records)
	imageUrls := f.getImageUrls(&imageResponseStruct)
	for index, url := range imageUrls {
		f.downloadImageUrl(url, tree, index)
	}

}

// Return a list of urls to download images from
func (f ForestryImageScraper) getImageUrls(imageResponseStruct *data.ImageResponse) []string {
	imageUrlList := make([]string, 0)
	for _, imageMap := range imageResponseStruct.Rows {
		imageNum := strconv.FormatInt(int64(imageMap[imagenum]), 10)
		formattedImageUrl := fmt.Sprintf(imageUrl, imageNum)
		imageUrlList = append(imageUrlList, formattedImageUrl)
	}
	return imageUrlList
}

// Downloads the image from the url
func (f ForestryImageScraper) downloadImageUrl(url, name string, index int) {
	c := colly.NewCollector()

	// "clean string" and create the directory to store the image
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "-")
	dir := fmt.Sprintf("%s/%s", f.downloadPath, name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory %s was not created, creating it!\n", dir)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Save the image
	c.OnResponse(func(r *colly.Response) {
		fileName := fmt.Sprintf("%s/%s/%d.jpg", f.downloadPath, name, index)
		err := r.Save(fileName)
		if err != nil {
			panic(err)
		}
	})
	c.Visit(url)
}

func (f ForestryImageScraper) ScrapeTree(in <-chan data.TreeJson, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		tree, ok := <-in
		if !ok {
			break
		}
		f.ScrapeImages(tree)
	}
}
