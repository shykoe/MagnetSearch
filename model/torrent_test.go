package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseTorrents(t *testing.T) {
	var r map[string]interface{}
	content, _ := ioutil.ReadFile("./test_es_content")
	json.Unmarshal(content, &r)
	res, _ :=  ParseTorrents(r)
	for _, item := range res{
		fmt.Println(item.Name)
	}
}