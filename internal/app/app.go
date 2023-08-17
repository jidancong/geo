package app

import (
	"fmt"
	"log"

	"github.com/jidancong/geo/config"
	"github.com/jidancong/geo/internal/usecase"
	"github.com/jidancong/geo/internal/usecase/repo"
	"github.com/jidancong/geo/internal/usecase/webapi"
	"github.com/jidancong/geo/pkg/es"
	"github.com/jidancong/geo/pkg/kafka"
)

func Run(cfg *config.Config) {

	kafkaClient, err := kafka.NewKafkaPkg(cfg.Kafka.Host, cfg.Kafka.Topic, cfg.Kafka.GroupId, cfg.Kafka.NumPartitions, cfg.Kafka.Replication)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - kafka.New: %w", err))
	}

	esClient, err := es.NewElasticsearchClient(cfg.Elasticsearch.Host, cfg.Elasticsearch.PreName)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - es.New: %w", err))
	}
	uc := usecase.NewCreateIndexUseCase(
		repo.NewRepoElasticsearch(esClient),
		repo.NewConsumerKafka(kafkaClient),
		webapi.NewLookupIp2LocationWebAPI(cfg.IP2location.Path),
		webapi.NewLookupChunZhen(cfg.ChunZhen.Path),
	)

	uc.Create()

}
