package scraper

import (
	"github.com/Dhsieh/tree_scraper/data"
)

type Scraper interface {
	ScrapeImages(input interface{})
	ScrapeKeyWordImages(keyword string)
	ScrapeTreeDatas(treeData []data.TreeJson)
}
