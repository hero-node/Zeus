package rpcclient

import (
	"testing"
)

func TestCallETH(t *testing.T) {
	Call_ETH("http://106.14.187.240:8545", "web3_clientVersion", []interface{}{})
}
