package heronode

import (
	"zeus/rpcclient"
	"zeus/utils/global"
)

func Call_ETH(method string, params []interface{}) (*rpcclient.ETHResp, error) {
	host := global.Ethhost()
	return rpcclient.Call_ETH(host, method, params)
}

func Call_QTUM(method string, params []interface{}) (*rpcclient.QTUMResp, error) {
	host := global.Qtumhost()
	return rpcclient.Call_QTUM(host, method, params)
}
