package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Dhsieh/tree_scraper/config"
	"github.com/Dhsieh/tree_scraper/data"
	"github.com/Dhsieh/tree_scraper/scraper/forestryimages"
	"github.com/spf13/viper"
)

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

	var workerGroup sync.WaitGroup

	if *info {
		fmt.Println("Downloaing all tree species information")
		forestryscraper.DownloadAllTreeSpecies()
	} else if *conf != "" {
		siteScraper := configuration.GetScraper("../downloads/tree_data.json")

		if configuration.Keyword == "" {
			treeJsonMap := config.Setup("../downloads/tree_data.json")
			in := make(chan data.TreeJson, configuration.NumRoutines*4)

			for i := 0; i < configuration.NumRoutines; i++ {
				workerGroup.Add(1)
				go siteScraper.ScrapeTree(in, &workerGroup)
			}

			var treeList []data.TreeJson
			emptyNameCounter := 0
			// Create treeList from map or config
			if *all {
				fmt.Println("Scraping all trees")
				for _, treeData := range treeJsonMap {
					if treeData.CommonName != "" {
						treeList = append(treeList, treeData)
					} else {
						emptyNameCounter += 1
					}

				}
			} else {
				for _, treeSpecies := range configuration.Species {
					if treeData, ok := treeJsonMap[treeSpecies]; !ok {
						fmt.Printf("Could not find tree species for common name: %s\n", treeSpecies)
					} else {
						treeList = append(treeList, treeData)
					}
				}
			}

			for _, tree := range treeList {
				in <- tree
			}

			close(in)
			workerGroup.Wait()

			fmt.Printf("emptyNameCounter is %d\n", emptyNameCounter)

		} else {
			siteScraper.ScrapeImages(configuration.Keyword)
		}

	} else {
		fmt.Errorf("Did not specify anything!")
	}
}
