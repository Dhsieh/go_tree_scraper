package data

import (
	"encoding/json"
)

// Final struct that will be returned to be used
// Contains mostly useful information pertaining to a tree species
type TreeResponseData struct {
	CategoryId     int32
	CommonName     string
	ScientificName string
	NumberOfImages int32
	TreeType       string
}

// Exposed struct that will converte the treeResponse (from url) to TreeResponseData
type TreeResponse struct {
	treeResponse treeResponse
	Data         []TreeResponseData
}

// Custom UnmarshalJSON function
func (data *TreeResponse) UnmarshalJSON(b []byte, treeType string) error {
	err := json.Unmarshal(b, &data.treeResponse)
	for _, arr := range data.treeResponse.Data {
		treeResponseData := &TreeResponseData{int32(arr[0].(float64)), arr[1].(string), arr[2].(string), int32(arr[3].(float64)), treeType}
		data.Data = append(data.Data, *treeResponseData)
	}

	return err
}

func (data *TreeResponse) Columns() []string {
	return data.treeResponse.Columns
}

func (data *TreeResponse) RecordsTotal() int32 {
	return data.treeResponse.RecordsTotal
}

func (data *TreeResponse) RecordsFiltered() int32 {
	return data.treeResponse.RecordsFiltered
}

func (data *TreeResponse) Draw() int32 {
	return data.treeResponse.Draw
}

func (data *TreeResponse) Append(t *TreeResponse) {
	data.Data = append(data.Data, t.Data...)
}

// This is the actual response from forestryimage for each tree
// This is not exposed as there are certain fields that are not used and are unnecessary
type treeResponse struct {
	Columns         []string
	RecordsTotal    int32
	Data            [][]interface{}
	RecordsFiltered int32
	Draw            int32
}
