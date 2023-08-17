package usecase

import (
	"github.com/jidancong/geo/entity"
	"github.com/jidancong/geo/pkg/kafka"
)

type (
	LookupIPCNWebAPI interface {
		GetFullName(ipaddress string) string
	}

	LookupIPWebAPI interface {
		GetAll(ipaddress string) (entity.IPRecord, error)
	}

	ElasticsearchRepo interface {
		Create(record string) error
		BatchCreate(record string) error
	}

	KafkaRepo interface {
		Consumer(processor kafka.Processor)
	}
)
