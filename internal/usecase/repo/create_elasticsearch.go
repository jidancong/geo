package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jidancong/geo/pkg/es"

	"github.com/sourcegraph/conc/pool"
)

type RepoElasticsearch struct {
	esClient      *es.ElasticsearchClient
	recordChannel chan string
	once          sync.Once
}

func NewRepoElasticsearch(client *es.ElasticsearchClient) *RepoElasticsearch {
	recordChannel := make(chan string, 1000)
	return &RepoElasticsearch{esClient: client, recordChannel: recordChannel}
}

func (es *RepoElasticsearch) Create(record string) error {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(record), &result)
	if err != nil {
		return fmt.Errorf("es create doc error:%w", err)
	}

	timestamp := result["@timestamp"]
	timestampVal, ok := timestamp.(string)
	if !ok {
		return fmt.Errorf("timestamp error:%w", err)
	}
	t, err := time.Parse(time.RFC3339, timestampVal)
	if err != nil {
		return fmt.Errorf("time parse error:%w", err)
	}
	y := t.Format("2006")
	m := t.Format("01")
	d := t.Format("02")
	// logstash-2023.08.10
	indexName := es.esClient.PreName + y + "." + m + "." + d

	_, err = es.esClient.ClientTypedClient.Index(indexName).Request(result).Do(context.TODO())
	return err
}

func (es *RepoElasticsearch) BatchCreate(record string) error {
	es.recordChannel <- record
	es.once.Do(func() {
		go func() {
			p := pool.New().WithMaxGoroutines(500)
			for elem := range es.recordChannel {
				elem := elem
				p.Go(func() {
					es.Create(elem)
				})
			}
			p.Wait()

		}()
	})
	return nil
}
