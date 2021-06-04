package scraper

import (
	"github.com/Dhsieh/tree_scraper/data"
)

type Scraper interface {
	ScrapeImages(treeData data.TreeJson)
	ScrapeAllTrees()
}
