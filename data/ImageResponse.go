package data

type ImageResponse struct {
	Rows    []map[string]int32
	Page    int32
	Records int32
	Total   int32
}
