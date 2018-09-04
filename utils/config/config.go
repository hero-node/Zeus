package config

import (
	"github.com/Mercy-Li/Goconfig/config"
)

func GetHttpPort() string {
	httpPort, err := config.GetConfigString("api.listen")
	if err != nil {
		return ":9198"
	}
	return httpPort
}

func GetEthHost() string {
	ethHost, err := config.GetConfigString("ethhost")
	if err != nil {
		return "http://localhost:8545"
	}
	return ethHost
}

func GetQtumHost() string {
	qtumHost, err := config.GetConfigString("qtumhost")
	if err != nil {
		return "http://localhost:3889"
	}
	return qtumHost
}

func GetIpfsHost() string {
	ipfsHost, err := config.GetConfigString("ipfshost")
	if err != nil {
		return "http://localhost:5001"
	}
	return ipfsHost
}

func GetQtumUserAndPassowrd() []string {
	user, err := config.GetConfigString("qtumuser")
	password, err := config.GetConfigString("qtumpassword")
	if err != nil {
		return []string{"123", "123"}
	}
	return []string{user, password}
}
