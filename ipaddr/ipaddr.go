package ipaddr

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ErrorParse = errors.New("host or ip parsing error")
var ErrorFile = errors.New("filename is nil")

//exclude net id list
var netids = []string{"/8", "/16", "/64", "128"}

func HostIPv4() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if !exclude(addr.String(), netids) {
			var ip = addr.String()
			i := strings.Index(ip, "/")
			if i > 0 {
				ip = ip[:i]
			}
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

func ParseFile(strFilename string) (hosts []string, err error) {

	if strFilename == "" {
		return nil, ErrorFile
	}

	var filehost []string
	filehost, _ = readFromFile(strFilename)
	hosts = append(hosts, filehost...)

	var result []string
	t := map[string]struct{}{}
	for _, host := range hosts {
		if _, ok := t[host]; !ok {
			t[host] = struct{}{}
			result = append(result, host)
		}
	}
	hosts = result
	if len(hosts) == 0 {
		err = ErrorParse
		return nil, err
	}
	return hosts, nil
}

func ParseIP(ip string) (hosts []string) {
	if strings.Contains(ip, ",") {
		IPList := strings.Split(ip, ",")
		var ips []string
		for _, v := range IPList {
			ips = parseIP(v)
			hosts = append(hosts, ips...)
		}
	} else {
		hosts = parseIP(ip)
	}
	return hosts
}

func exclude(a string, ss []string) bool {
	for _, id := range ss {
		if strings.Contains(a, id) {
			return true
		}
	}
	return false
}

func parseIP(ip string) []string {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case strings.HasSuffix(ip, "/8"): //net id 8
		return parseNetID8(ip)
	case strings.Contains(ip, "/"): //net id 16/24/...
		return convert2Range(ip)
	case strings.Contains(ip, "-"): //ip range split by -
		return parseIPRange(ip)
	case reg.MatchString(ip): //domain
		return []string{ip}
	default: //single ip
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}
}

func parseIPRange(ip string) []string {
	iprange := strings.Split(ip, "-")
	testIP := net.ParseIP(iprange[0])
	var ips []string
	if len(iprange[1]) < 4 {
		Range, err := strconv.Atoi(iprange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(iprange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(iprange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			ips = append(ips, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		ip1 := strings.Split(iprange[0], ".")
		ip2 := strings.Split(iprange[1], ".")
		if len(ip1) != 4 || len(ip2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(ip1[i])
			ip2, err2 := strconv.Atoi(ip2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			ips = append(ips, ip)
		}
	}
	return ips
}

func convert2Range(host string) (hosts []string) {
	_, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		return
	}
	hosts = parseIPRange(ipRange(ipNet))
	return
}

// get start and end ip
func ipRange(c *net.IPNet) string {
	start := c.IP.String()
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	end := bcst.String()
	return fmt.Sprintf("%s-%s", start, end)
}

//逐行读取IP地址信息
func readFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open %s error, %v", filename, err)
		os.Exit(0)
	}
	defer file.Close()
	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			host := ParseIP(text)
			content = append(content, host...)
		}
	}
	return content, nil
}

func parseNetID8(ip string) []string {
	realIP := ip[:len(ip)-2]
	testIP := net.ParseIP(realIP)

	if testIP == nil {
		return nil
	}

	iprange := strings.Split(ip, ".")[0]
	var ips []string
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, 1))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, 2))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, 4))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, 5))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, randn(6, 55)))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, randn(56, 100)))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, randn(101, 150)))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, randn(151, 200)))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, randn(201, 253)))
			ips = append(ips, fmt.Sprintf("%s.%d.%d.%d", iprange, a, b, 254))
		}
	}
	return ips
}

func randn(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
