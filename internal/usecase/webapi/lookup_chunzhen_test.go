package webapi

import (
	"testing"
)

func TestChunZhen(t *testing.T) {
	q := NewLookupChunZhen("qqwry.dat")
	q.GetFullName("223.5.5.5")

}
