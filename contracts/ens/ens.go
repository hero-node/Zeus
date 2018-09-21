package ens

import (
	"log"
	"zeus/utils/config"
	"zeus/utils/global"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/contracts/ens/contract"
)

func ensContract() *contract.ENS, err {
	ethHost := config.GetValidEthHost()
	conn, err := ethclient.Dial(ethHost)
	if err != nil {
		return nil, err
	}
	
	var ensAddr := common.HexToAddress("0x314159265dD8dbb310642f98f50C066173C1259b")
	if global.TEST_NET {
		ensAddr = common.HexToAddress("0x112234455c3a32fd11230c42e7bccd4a84e02010")
	}
	
	ens, err := contract.NewENS(ensAddr, conn)
	if err != nil {
		return nil, err
	}
	return ens, nil
}

func ensParentNode(name string) (common.Hash, common.Hash) {
	parts := strings.SplitN(name, ".", 2)
	label := crypto.Keccak256Hash([]byte(parts[0]))
	if len(parts) == 1 {
		return [32]byte{}, label
	} else {
		parentNode, parentLabel := ensParentNode(parts[1])
		return crypto.Keccak256Hash(parentNode[:], parentLabel[:]), label
	}
}

func ensNode(name string) common.Hash {
	parentNode, parentLabel := ensParentNode(name)
	return crypto.Keccak256Hash(parentNode[:], parentLabel[:])
}

func EnsToAddr(name string) string {
	ens, err := ensConteact()
	if err != nil {
		log.Printf("Get Ens contract failed: %v", err)
		return ""
	}
	namehash := ensNode(name)
	resolverAddr, err := ens.Resolver(nil, namehash)
	if err != nil {
		log.Printf("Get resolver failed: %v", err)
		return ""
	}
	addr, err := resolverContract.Addr(nil, namehash)
	if err != nil {
		log.Printf("Resolver get addr failed: %v", err)
		return ""
	}
	
	return addr.Hex()
}

func EnsToIpfs(name string) string {
	ens, err := ensConteact()
	if err != nil {
		log.Printf("Get Ens contract failed: %v", err)
		return ""
	}
	namehash := ensNode(name)
	resolverAddr, err := ens.Resolver(nil, namehash)
	if err != nil {
		log.Printf("Get resolver failed: %v", err)
		return ""
	}
	addr, err := resolverContract.Content(nil, namehash)	
	
}
