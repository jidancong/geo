package webapi

import (
	"fmt"
	"log"

	"github.com/jidancong/geo/entity"

	"github.com/ip2location/ip2location-go/v9"
)

type LookupIp2LocationWebAPI struct {
	db *ip2location.DB
}

func NewLookupIp2LocationWebAPI(path string) *LookupIp2LocationWebAPI {
	db, err := ip2location.OpenDB(path)
	if err != nil {
		log.Fatalf("new ip2location error:%s", err)
	}
	return &LookupIp2LocationWebAPI{db}
}

func (l *LookupIp2LocationWebAPI) GetAll(ipaddress string) (entity.IPRecord, error) {
	result, err := l.db.Get_all(ipaddress)
	if err != nil {
		return entity.IPRecord{}, fmt.Errorf("get ip record error: %s", err)
	}

	return entity.IPRecord{
		Country:   result.Country_long,
		Region:    result.Region,
		City:      result.City,
		Latitude:  float64(result.Latitude),
		Longitude: float64(result.Longitude),
	}, nil
}
