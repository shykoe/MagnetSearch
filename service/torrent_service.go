package service

import (
	"MagnetSearch/dao"
	"MagnetSearch/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func BuildSearchService(dao *dao.TorrentDaoImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchContent := c.Query("content")
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		if size <= 0 || size >= 50{
			size = 50
		}
		pos, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		torrents, count, err := dao.FindTorrents(searchContent, int32(pos), int32(size))
		if err != nil{
			c.JSON(500, err.Error())
			return
		}
		var result model.SearchResult
		result.PageSize = size
		result.Count = int(count)
		result.Torrent = torrents
		// build link
		var urls []string
		for _, item := range torrents {
			url := fmt.Sprintf("magnet:?xt=urn:btih:%s", item.InfoHash)
			urls = append(urls, url)
		}
		result.Links = urls
		c.JSON(http.StatusOK,result)
	}
}
