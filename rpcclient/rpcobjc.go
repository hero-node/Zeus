package rpcclient

type ETHResp struct {
	ID      string      `json:"id"`
	JSONPRC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type QTUMResp struct {
	ID     string      `json:"id"`
	error  string      `json:"error"`
	Result interface{} `json:"result"`
}
