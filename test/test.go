package main

import (
	"fmt"
	"log"
	//	"log"
	"strings"

	//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens/contract"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// 6241b7c237fe1fefa89b9e99580159de2621ee4a679d27afbb6289fcd74f09d1
	//	key, _ := crypto.HexToECDSA("6241b7c237fe1fefa89b9e99580159de2621ee4a679d27afbb6289fcd74f09d1")

	conn, err := ethclient.Dial("https://mainnet.infura.io/v3/719be1b239a24d1e87a2e326be6c4384")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	//	auth := bind.NewKeyedTransactor(key)
	addr := common.HexToAddress("0x314159265dD8dbb310642f98f50C066173C1259b")
	//	ens, err := ens.NewENS(transactOpt, addr, conn)
	//	if err != nil {
	//		log.Fatalf("Failed to instantiate a Token contract: %v, %v", err, addr)
	//	}
	//	hash, err := ens.Resolve("heronode")
	//	if err != nil {
	//		log.Fatalf("Failed to resolve: %v", err)
	//	}
	//	fmt.Println(" ico hash:", hash.Hex())

	namehash := ensNode("heronode.eth")

	ens, err := contract.NewENS(addr, conn)
	if err != nil {
		log.Fatalln("Failed to new ens")
	}

	resolverAddr, err := ens.Resolver(nil, namehash)
	if err != nil {
		log.Fatalln("reslover failed")
	}

	resolverContract, _ := contract.NewPublicResolver(resolverAddr, conn)
	target, err := resolverContract.Addr(nil, namehash)
	fmt.Println(target.Hex())

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
