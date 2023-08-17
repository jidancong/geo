package usecase

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
)

type CreateIndexUseCase struct {
	repo  ElasticsearchRepo
	queue KafkaRepo
	api   LookupIPWebAPI
	cnapi LookupIPCNWebAPI
}

func NewCreateIndexUseCase(repo ElasticsearchRepo, queue KafkaRepo, api LookupIPWebAPI, cnapi LookupIPCNWebAPI) *CreateIndexUseCase {
	return &CreateIndexUseCase{repo, queue, api, cnapi}
}

func (c *CreateIndexUseCase) Create() {
	c.queue.Consumer(func(ctx context.Context, msg []byte) error {
		var result map[string]interface{}
		json.Unmarshal(msg, &result)

		ip, ok := result["client"].(string)
		if !ok {
			log.Println("String type conversion error")
			return nil
		}

		// 验证ip
		if ok := isIpV4(ip); !ok {
			log.Println("IP address validation failed")
			return nil
		}

		// 纯真数据查询
		result["fullname"] = c.cnapi.GetFullName(ip)

		// ip2location数据库查询
		record, _ := c.api.GetAll(ip)
		result["latitude"] = record.Latitude
		result["longitude"] = record.Longitude
		result["country"] = record.Country
		result["region"] = record.Region
		result["city"] = record.City

		b, _ := json.Marshal(result)
		err := c.repo.BatchCreate(string(b))
		// err := c.repo.Create(string(b))
		log.Println(result)
		return err
	})
}

func isIpV4(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return strings.Contains(ipstr, ".")
}
