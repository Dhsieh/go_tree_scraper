package config

import (
	"fmt"
	"github.com/Dhsieh/tree_scraper/scraper"
	"github.com/Dhsieh/tree_scraper/scraper/bing"
	"github.com/Dhsieh/tree_scraper/scraper/forestryimages"
)

type Configuration struct {
	Site         string
	Species      []string
	DownloadPath string
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
		return bingscraper.CreateScraper(c.getFullPath(), Setup(jsonPath))
	}
}
