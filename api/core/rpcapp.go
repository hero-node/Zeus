package core

import (
	"zeus/rpcclient"
	"zeus/utils/config"
)

func Call_ETH(method string, params []interface{}) (*rpcclient.ETHResp, error) {
	host := config.GetEthHost()
	return rpcclient.Call_ETH(host, method, params)
}

func Call_QTUM(method string, params []interface{}) (*rpcclient.QTUMResp, error) {
	host := config.GetQtumHost()
	return rpcclient.Call_QTUM(host, method, params)
}
