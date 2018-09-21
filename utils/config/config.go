package config

import (
	"context"
	"log"
	"time"
	"zeus/utils/global"

	"github.com/Mercy-Li/Goconfig/config"
	"github.com/ethereum/go-ethereum/ethclient"
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

var validedEthHost = ""
var lastTime time.Time

func GetValidEthHost() string {
	var validEthHost string
	if global.TEST_NET {
		validEthHost = "https://ropsten.infura.io/v3/719be1b239a24d1e87a2e326be6c4384"
	} else {
		validEthHost = "https://mainnet.infura.io/v3/719be1b239a24d1e87a2e326be6c4384"
	}

	duration := time.Since(lastTime).Hours()
	if validedEthHost == "" || duration > 1 {
		lastTime = time.Now()
		localHost := GetEthHost()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(localHost)
		if err != nil {
			validedEthHost = validEthHost
			return validedEthHost
		}
		block_local, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			validedEthHost = validEthHost
			return validedEthHost
		}

		client, err = ethclient.Dial(validEthHost)
		if err != nil {
			log.Fatal("No valid node available")
		}
		block_valid, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			log.Fatal("No valid node available")
		}
		if int(block_local.NumberU64()-block_valid.NumberU64()) > -5 {
			validedEthHost = localHost
			return validedEthHost
		}
	}

	return validedEthHost
}
