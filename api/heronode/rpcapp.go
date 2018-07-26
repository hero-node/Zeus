package heronode

import (
	"zeus/rpcclient"

	"github.com/Mercy-Li/Goconfig/config"
)

func Call_ETH(method string, params []interface{}) (*rpcclient.ETHResp, error) {
	host, err := config.GetConfigString("ethhost")
	if err != nil {
		return nil, err
	}
	return rpcclient.Call_ETH(host, method, params)
}
