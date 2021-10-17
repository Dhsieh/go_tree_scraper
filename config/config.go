package config

import (
	"fmt"
	"github.com/Dhsieh/tree_scraper/scraper"
	"github.com/Dhsieh/tree_scraper/scraper/bing"
	"github.com/Dhsieh/tree_scraper/scraper/forestryimages"
)

// All configuration options
// Site: 			string      which website to scrape, creates a specific scraper for that website
// Species: 		[]string    list of a species' common names to scrape from the Site
// DownloadPath: 	String  	folder to scrape and download the images to
// Images: 			int 		number of images to scrape per species
// NumRoutines: 	int 	 	number of Go routines to create
type Configuration struct {
	Site         string
	Keyword      string
	Species      []string
	DownloadPath string
	Images       int
	NumRoutines  int
}

func (c Configuration) getFullPath() string {
	return fmt.Sprintf("%s/%s", c.DownloadPath, c.Site)
}

func (c Configuration) GetScraper(jsonPath string) scraper.Scraper {
	if c.Site == "bugwood" {
		fmt.Println("Scraping bugwood")
		return forestryscraper.CreateScraper(c.getFullPath(), Setup(jsonPath))
	} else {
		fmt.Println("Scraping bing")
		bing := bingscraper.CreateScraper(c.getFullPath(), Setup(jsonPath), c.Images, c.NumRoutines)
		return bing
	}
}

func (c Configuration) CreateBingScraper() bingscraper.BingScraper {
	return bingscraper.CreateScraper(c.DownloadPath, nil, c.Images, c.NumRoutines)
}

func (c Configuration) String() string {
	return fmt.Sprintf("Site: %s\nKeyword: %s\nSpecies: %s\nDownloadPath: %s\nImages: %d\nNumRoutines: %d\n", c.Site, c.Keyword, c.Species, c.DownloadPath, c.Images, c.NumRoutines)
}
