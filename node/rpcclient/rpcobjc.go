package rpcclient

type ETHResp struct {
	ID      int    `json:"id"`
	JSONPRC string `json:"jsonrpc"`
	Result  string `json:"result"`
}
