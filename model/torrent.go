package model

import "time"

type TorrentFile struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
}
type Torrent struct {
	InfoHash string         `json:"Hash"`
	Name     string         `json:"Name"`
	Length   int64          `json:"Length"`
	DiscoverFrom string `json:"DiscoverFrom"`
	DiscoverTime time.Time `json:"DiscoverTime"`
	Files    []*TorrentFile `json:"Files"`
}
func ParseTorrents(content map[string]interface{})([]*Torrent, error){
	hits := content["hits"].(map[string]interface{})["hits"].([]interface{})
	var result []*Torrent
	for _, hit := range hits{
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		instance := new(Torrent)
		instance.InfoHash = source["Hash"].(string)
		instance.DiscoverFrom = source["DiscoverFrom"].(string)
		instance.Name = source["Name"].(string)
		instance.DiscoverTime, _ = time.Parse("2006-01-02T15:04:05Z0700", source["DiscoverTime"].(string))
		instance.Length = int64(source["Length"].(float64))
		var files []*TorrentFile
		for _, item := range source["Files"].([]interface{}){
			file := new(TorrentFile)
			file.Name = item.(map[string]interface{})["Name"].(string)
			file.Size = int64(item.(map[string]interface{})["Size"].(float64))
			files = append(files, file)
		}
		instance.Files = files
		result = append(result, instance)
	}
	return result,nil
}