package store

type KNNSearchResponse struct {
	Key      string
	Distance float32
}

type SearchResponse struct {
	Raw    string
	Answer string
	Vector []float32
}

type SearchResponsePair struct {
	Key      string
	Distance float32
	Vector   []float32
	Raw      string
	Answer   string
}
