package ipaddr

import (
	"fmt"
	"testing"
)

func TestParseIP(t *testing.T) {
	var ip = "192.168.1.0/24"
	list := ParseIP(ip)
	for _, v := range list {
		fmt.Printf("ip %s \n", v)
	}
	fmt.Printf("total [%d] ip address", len(list))
}

func TestParseFile(t *testing.T) {

}
