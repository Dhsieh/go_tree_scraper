package data

// Information gained for a specific tree spedcies from foresty images.
type TreeJson struct {
	CategoryId     int32
	CommonName     string
	ScientificName string
	NumberOfImages int32
}

// List of TreeJson that is obtained from foresty image url
type TreeJsons struct {
	TreeJsons []TreeJson `json:"Data"`
}
