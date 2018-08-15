package global

import (
	"log"

	"github.com/Mercy-Li/Goconfig/config"
)

const (
	BTC  = "btc"
	ETH  = "eth"
	QTUM = "qtum"
)

func Ethhost() string {
	ethhost, err := config.GetConfigString("ethhost")
	if err != nil {
		log.Fatal(err)
	}
	return ethhost
}

func IpfsHost() string {
	ipfshost, err := config.GetConfigString("ipfshost")
	if err != nil {
		log.Fatal(err)
	}
	return ipfshost
}
