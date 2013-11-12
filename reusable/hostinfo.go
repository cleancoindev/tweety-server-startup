package reusable

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func filter(s []string, fn func(string) bool) []string {
	var p []string // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func is_ipv4(in string) bool {
	return !strings.Contains(in, ":")
}

func get_public_ip() string {
	resp, err := http.Get("http://api.externalip.net/ip")
	if err != nil {
		log.Fatal(err)
	}

	return response_to_string(resp)
}

func GetHostnameAndIps() (string, string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	local_ip_addresses, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal(err)
	}

	public_ip := get_public_ip()

	ip_addresses := fmt.Sprintf("%s, %s", public_ip, strings.Join(filter(local_ip_addresses, is_ipv4), ", "))

	return hostname, ip_addresses
}