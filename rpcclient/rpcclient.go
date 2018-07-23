package rpcclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Call_ETH(host string, method string, params []interface{}) (interface{}, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "go-heronode",
	}
	data["method"] = method
	data["params"] = params

	jsonstring, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(jsonstring))
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body", string(body))

	return body, nil
}
