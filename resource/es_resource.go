package resource

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

var esConf *elasticsearch7.Config

//初始化配置
func InitConf(address string, user string, password string) {
	esConf = &elasticsearch7.Config{
		Addresses: []string{
			address,
		},
		Username: user,
		Password: password,
	}
}
func GetESClient()(*elasticsearch7.Client,error)  {

	client, err := elasticsearch7.NewClient(*esConf)
	if err !=nil{
		return nil,err
	}
	return client,nil
}