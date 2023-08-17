package es

import (
	"fmt"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	ClientTypedClient *elasticsearch8.TypedClient
	PreName           string
}

func NewElasticsearchClient(host string, preName string) (*ElasticsearchClient, error) {
	if len(preName) <= 0 {
		preName = "logstash-"
	}

	cfg := elasticsearch8.Config{
		Addresses: []string{host},
	}
	es8, err := elasticsearch8.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("new elasticsearch error:%s", err)
	}
	return &ElasticsearchClient{es8, preName}, nil
}
