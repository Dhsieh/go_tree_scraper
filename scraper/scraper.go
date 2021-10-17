package scraper

import (
	"github.com/Dhsieh/tree_scraper/data"
)

type Scraper interface {
	ScrapeImages(input interface{})
	ScrapeKeywordImages(keyword string)
	ScrapeTreeDatas(treeData []data.TreeJson)
}
