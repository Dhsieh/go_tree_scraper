package data

// Struct that contains information from a forestryimage url on a specific species
type ImageResponse struct {
	Rows    []map[string]int32
	Page    int32
	Records int32
	Total   int32
}
