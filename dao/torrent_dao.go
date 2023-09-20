package dao

import (
	"MagnetSearch/model"
	"bytes"
	"context"
	"encoding/json"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	log "github.com/sirupsen/logrus"
	"strings"
)

type TorrentDao interface {
	FindTorrents(content string, CurrentPage int32, pageSize int32) ([]*model.Torrent, int64, error)
}
type TorrentDaoImpl struct {
	Es    *elasticsearch7.Client
	Index string
}

func (t TorrentDaoImpl) FindTorrents(content string, CurrentPage int32, pageSize int32) ([]*model.Torrent, int64, error) {
	var (
		r map[string]interface{}
	)
	// build search query
	var buf bytes.Buffer
	//
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"type":    "best_fields",
							"lenient": true,
							"query":   content,
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"DiscoverTime": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"size":             pageSize,
		"from":             (CurrentPage - 1) * pageSize,
		"track_total_hits": "true",
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Errorf("Error encoding query: %s", err)
	}
	res, err := t.Es.Search(
		t.Es.Search.WithContext(context.Background()),
		t.Es.Search.WithIndex(t.Index),
		t.Es.Search.WithBody(&buf),
		t.Es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
		return nil, 0, err
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Errorf("Error parsing the response body: %s", err)
		return nil, 0, err
	}
	log.Infof("res : %v+", r)
	size := r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	torrents, err := model.ParseTorrents(r)
	if err != nil {
		return nil, 0, err
	}
	return torrents, int64(size), nil
}

func (t TorrentDaoImpl) InsertTorrent2Es(torrent *model.Torrent) error {
	if t.Es == nil {
		log.Error("es nil ")
		return nil
	}
	bytes, err := json.Marshal(torrent)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	req := esapi.IndexRequest{
		Index:      t.Index,
		DocumentID: torrent.InfoHash,
		Body:       strings.NewReader(string(bytes)),
		Refresh:    "false",
	}
	res, err := req.Do(context.Background(), t.Es)
	if err != nil {
		log.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Errorf("[%s] Error indexing document ID=%d msg=%s", res.Status(), torrent.InfoHash, res.String())
	}
	return nil
}
