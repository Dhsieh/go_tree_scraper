package main

import (
	"flag"
	"fmt"

	"github.com/Dhsieh/tree_scraper/config"
	"github.com/Dhsieh/tree_scraper/scraper/forestryimages"
	"github.com/spf13/viper"
)

func main() {
	all := flag.Bool("all", false, "download all tree species images or not")
	info := flag.Bool("info", false, "Download all possible tree species into a json file")
	conf := flag.Bool("conf", false, "Config to use")

	flag.Parse()

	viper.SetConfigName("confg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	var configuration config.Configuration

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(*all)
	viper.Unmarshal(&configuration)

	if *info {
		fmt.Println("Downloaing all tree species information")
		forestryscraper.DownloadAllTreeSpecies()
	} else if *conf {
		treeJsonMap := config.Setup("../downloads/tree_data.json")
		siteScraper := configuration.GetScraper("../downloads/tree_data.json")
		if *all {
			fmt.Println("Scraping all trees")
			siteScraper.ScrapeAllTrees()
		} else {
			for _, treeSpecies := range configuration.Species {
				fmt.Println(treeSpecies)
				if treeData, ok := treeJsonMap[treeSpecies]; !ok {
					fmt.Printf("Could not find tree species for common name: %s\n", treeSpecies)
				} else {
					siteScraper.ScrapeImages(treeData)
				}
			}
		}

	} else {
		fmt.Errorf("Did not specify anything!")
	}

}
