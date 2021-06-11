package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
	hybridCounter := 0
	for _, treeJson := range treeJsons.TreeJsons {
		if strings.HasPrefix(treeJson.ScientificName, "x ") || strings.Contains(treeJson.ScientificName, " x ") {
			hybridCounter += 1
			fmt.Printf("Igoring tree species %s as it is a hybrid!\n", treeJson.ScientificName)
		} else {
			treeListMap[treeJson.ScientificName] = treeJson
		}
	}

	fmt.Printf("Total of %d hybrids are ignored", hybridCounter)
	return treeListMap
}
