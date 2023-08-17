package webapi

import (
	"log"

	"github.com/xiaoqidun/qqwry"
)

type LookupChunZhen struct {
}

func NewLookupChunZhen(path string) *LookupChunZhen {
	if err := qqwry.LoadFile(path); err != nil {
		panic(err)
	}
	return &LookupChunZhen{}
}

func (l *LookupChunZhen) GetFullName(ipaddress string) string {
	city, isp, err := qqwry.QueryIP(ipaddress)
	if err != nil {
		log.Println(err)
		return ""
	}
	return city + "-" + isp
}
