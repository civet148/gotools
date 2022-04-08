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
	fmt.Printf("PARSE total [%d] ip address\n\n", len(list))
}

func TestParseFile(t *testing.T) {

}

func TestHostIPv4(t *testing.T) {
	list, _ := HostIPv4()
	for _, v := range list {
		fmt.Printf("ip %s \n", v)
	}
	fmt.Printf("HOST total [%d] ip address\n\n", len(list))
}

func TestNetIPv4(t *testing.T) {
	list, _ := NetIPv4()
	for _, v := range list {
		fmt.Printf("ip %s \n", v)
	}
	fmt.Printf("NET total [%d] ip address\n\n", len(list))
}
