package main

import (
	"flag"
	"fmt"

	"github.com/Dhsieh/tree_scraper/config"
	"github.com/Dhsieh/tree_scraper/data"
	"github.com/Dhsieh/tree_scraper/scraper/forestryimages"
	"github.com/spf13/viper"
)

// Create Tree Slice depending on config.
func createTreeSlice(all bool, treeJsonMap map[string]data.TreeJson, configuration config.Configuration) []data.TreeJson {
	var treeSlice []data.TreeJson
	counter := 0
	if all {
		for _, treeData := range treeJsonMap {
			if treeData.CommonName != "" {
				treeSlice = append(treeSlice, treeData)
			} else {
				counter++
			}
		}
		fmt.Printf("There are %d empty common names", counter)
	} else {
		for _, treeSpecies := range configuration.Species {
			if treeData, ok := treeJsonMap[treeSpecies]; !ok {
				counter++
			} else {
				treeSlice = append(treeSlice, treeData)
			}
		}
	}

	return treeSlice
}

func main() {
	// all downloads all the tree species in tree_data.json
	// info creates tree_data.json
	// conf contains certain configuration options as well as giving the option to scrape only a certain number of trees
	all := flag.Bool("all", false, "download all tree species images or not")
	info := flag.Bool("info", false, "Download all possible tree species into a json file")
	conf := flag.String("conf", "", "Config to use")

	flag.Parse()

	viper.SetConfigName(*conf)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	var configuration config.Configuration

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&configuration)

	if *info {
		fmt.Println("Downloaing all tree species information")
		forestryscraper.DownloadAllTreeSpecies()
	} else if *conf != "" {
		siteScraper := configuration.GetScraper("../downloads/tree_data.json")

		if configuration.Keyword == "" && len(configuration.Species) > 0 {
			treeJsonMap := config.Setup("../downloads/tree_data.json")
			treeSlice := createTreeSlice(*all, treeJsonMap, configuration)
			siteScraper.ScrapeImages(treeSlice)
		} else if len(configuration.Species) == 0 && configuration.Keyword != "" {
			siteScraper.ScrapeImages(configuration.Keyword)
		} else {
			panic(fmt.Errorf("Did not specify a keyword or a list of tree species!"))
		}

	} else {
		panic(fmt.Errorf("Did not specify anything!"))
	}
}
