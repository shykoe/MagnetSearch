package model

type SearchResult struct {
	Count       int      `json:"Count"`
	PageSize    int      `json:"PageSize"`
	CurrentPage int      `json:"CurrentPage"`
	Torrent     []*Torrent `json:"Torrent"`
	Links       []string  `json:"Links"`
}
