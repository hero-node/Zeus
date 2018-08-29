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

func ApiListenPort() string {
	port, err := config.GetConfigString("api.listen")
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func Qtumhost() string {
	qtumhost, err := config.GetConfigString("qtumhost")
	if err != nil {
		log.Fatal(err)
	}
	return qtumhost
}

func QtumUserAndPassword() []string {
	user, err := config.GetConfigString("qtumuser")
	password, err := config.GetConfigString("qtumpassword")
	if err != nil {
		return []string{"", ""}
	}
	return []string{user, password}
}
