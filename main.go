package main

import (
"MagnetSearch/dao"
"MagnetSearch/resource"
"MagnetSearch/service"
"fmt"
"github.com/gin-gonic/gin"
log "github.com/sirupsen/logrus"
"gopkg.in/yaml.v2"
"io/ioutil"
)

func main() {
	config := make(map[string]string)
	data, err := ioutil.ReadFile("./config.yml")
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	resource.InitConf(config["es_addr"], config["es_user"], config["es_pswd"])
	var torrentDao dao.TorrentDaoImpl
	torrentDao.Es, _ = resource.GetESClient()
	torrentDao.Index = config["es_index"]
	var router = gin.Default()
	router.GET("/torrent/search", service.BuildSearchService(&torrentDao))
	router.Run(fmt.Sprintf(":%s", config["port"]))
}
