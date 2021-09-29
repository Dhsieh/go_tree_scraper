package scraper

import (
	"sync"

	"github.com/Dhsieh/tree_scraper/data"
)

type Scraper interface {
	ScrapeImages(input interface{})
	ScrapeKeyWordImages(keyword string)
	ScrapeTreeData(treeData data.TreeJson)
	ScrapeTree(in <-chan data.TreeJson, wg *sync.WaitGroup)
}
