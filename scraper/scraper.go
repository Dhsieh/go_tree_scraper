package scraper

import (
	"sync"

	"github.com/Dhsieh/tree_scraper/data"
)

type Scraper interface {
	ScrapeImages(treeData data.TreeJson)
	ScrapeTree(in <-chan data.TreeJson, wg *sync.WaitGroup)
}
