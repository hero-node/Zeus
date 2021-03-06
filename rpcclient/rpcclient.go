package rpcclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"zeus/utils/config"
)

func Call_ETH(host string, method string, params []interface{}) (*ETHResp, error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObjc := new(ETHResp)
	err = json.Unmarshal(body, &respObjc)
	if err != nil {
		return nil, err
	}

	return respObjc, nil
}

func Call_QTUM(host string, method string, params []interface{}) (*QTUMResp, error) {
	data := map[string]interface{}{
		"jsonrpc": "1.0",
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

	req.Header.Set("content-type", "text/plain")
	up := config.GetQtumUserAndPassowrd()
	req.SetBasicAuth(up[0], up[1])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObjc := new(QTUMResp)
	err = json.Unmarshal(body, &respObjc)
	if err != nil {
		return nil, err
	}

	return respObjc, nil
}
