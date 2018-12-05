package hns

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"zeus/nameservice/hns/contracts"
	"zeus/utils/config"
	"zeus/utils/global"
)

var hnsContract *contracts.HNS
var initialized bool

var rosptenAddress = "0x68d2e1a0aaebd98ef4637aa2941921f050f1052d"
var mainAddress = ""

func newHnsContract() (*contracts.HNS, error) {
	if !initialized {
		initialized = true
		contractAddress := mainAddress
		if global.TEST_NET {
			contractAddress = rosptenAddress
		}
		ethHost := config.GetValidEthHost()
		client, _ := ethclient.Dial(ethHost)
		hnsContract, err := contracts.NewHNS(common.HexToAddress(contractAddress), client)
		if err != nil {
			return nil, err
		}
	}
	return hnsContract, nil
}

func SetHNS(name string, key string) (error) {

}

func GetHNS(key string) (name string, err error) {

}

