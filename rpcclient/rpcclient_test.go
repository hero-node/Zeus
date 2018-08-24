package rpcclient

import (
	"fmt"
	"testing"
)

func TestCallETH(t *testing.T) {
	resp, err := Call_ETH("http://106.14.187.240:8545", "web3_clientVersion", []interface{}{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestCallQTUM(t *testing.T) {
	resp, err := Call_QTUM("http://106.14.187.240:3889", "getinfo", []interface{}{})
	fmt.Println(resp)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
