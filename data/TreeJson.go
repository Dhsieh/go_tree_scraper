package data

type TreeJson struct {
	CategoryId     int32
	CommonName     string
	ScientificName string
	NumberOfImages int32
}

type TreeJsons struct {
	TreeJsons []TreeJson `json:"Data"`
}
