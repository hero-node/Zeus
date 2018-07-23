package rpcclient

import (
	"net/http"
)

func Call(host string, method string, parameters interface{}) map[string]interface{} {
	client := &http.Client{}
	req, err := http.NewRequest("POST", host)
	req.Header.Set("content-type": "text/plain")
}

