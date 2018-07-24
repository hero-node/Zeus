package rpcclient

type ETHResp struct {
	ID      string `json:"id"`
	JSONPRC string `json:"jsonrpc"`
	Result  string `json:"result"`
}
