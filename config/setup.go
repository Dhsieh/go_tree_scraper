package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Dhsieh/tree_scraper/data"
)

// Creates the map of all tree species from a file
func Setup(filename string) map[string]data.TreeJson {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var treeJsons data.TreeJsons
	err = json.Unmarshal(dat, &treeJsons)
	if err != nil {
		panic(err)
	}

	treeListMap := make(map[string]data.TreeJson)
	for _, treeJson := range treeJsons.TreeJsons {
		treeListMap[treeJson.CommonName] = treeJson
	}

	return treeListMap
}
