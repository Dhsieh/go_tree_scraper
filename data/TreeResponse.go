package data

import (
	"encoding/json"
)

type TreeResponseData struct {
	CategoryId     int32
	CommonName     string
	ScientificName string
	NumberOfImages int32
	TreeType       string
}

type TreeResponse struct {
	treeResponse treeResponse
	Data         []TreeResponseData
}

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

type treeResponse struct {
	Columns         []string
	RecordsTotal    int32
	Data            [][]interface{}
	RecordsFiltered int32
	Draw            int32
}
