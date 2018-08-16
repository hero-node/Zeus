package utils

import (
	"strings"
)

func ExtractIPv4(url string) string {
	if strings.Contains(url, "ip4") {
		split := strings.Split(url, "/")
		ip4 := split[2]
		return ip4
	}
	return ""
}

func ConstructUrl(host string) string {
	schema := "http"
	port := "9198"

	return schema + "://" + host + ":" + port
}
