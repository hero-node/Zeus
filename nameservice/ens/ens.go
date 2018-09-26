package ens

import (
	"encoding/hex"
	"log"
	"strings"
	"zeus/utils/config"
	"zeus/utils/global"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens/contract"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/multiformats/go-multihash"
)

func ensContract() (*contract.ENS, *ethclient.Client, error) {
	ethHost := config.GetValidEthHost()
	conn, err := ethclient.Dial(ethHost)
	if err != nil {
		return nil, nil, err
	}
	ensAddr := common.HexToAddress("0x314159265dD8dbb310642f98f50C066173C1259b")
	if global.TEST_NET {
		ensAddr = common.HexToAddress("0x112234455c3a32fd11230c42e7bccd4a84e02010")
	}

	ens, err := contract.NewENS(ensAddr, conn)
	if err != nil {
		return nil, nil, err
	}
	return ens, conn, nil
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
	ens, conn, err := ensContract()
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
	resolverContract, err := contract.NewPublicResolver(resolverAddr, conn)
	if err != nil {
		log.Printf("Get resolver contract failed: %v", err)
		return ""
	}
	addr, err := resolverContract.Addr(nil, namehash)
	if err != nil {
		log.Printf("Resolver get addr failed: %v", err)
		return ""
	}

	return addr.Hex()
}

// ======================== ipfs part ==============================
func EnsToContent(name string) string {
	ens, conn, err := ensContract()
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
	resolverContract, err := contract.NewPublicResolver(resolverAddr, conn)
	if err != nil {
		log.Printf("Get resolver contract failed: %v", err)
		return ""
	}

	content, err := resolverContract.Content(nil, namehash)
	if err != nil {
		log.Printf("Get resolver content failed: %v", err)
		return ""
	}
	return string(content[:])
}

func IpfsEncode(ipfs string) (string, error) {
	multi, err := multihash.FromB58String(ipfs)
	if err != nil {
		return "", err
	}

	buf := multi.HexString()
	return "0x" + buf, nil
}

func IpfsDecode(hash string) (string, error) {
	_hash := hash
	if hash[:2] == "0x" {
		_hash = hash[2:]
	}

	buf, _ := hex.DecodeString(_hash)
	mul := multihash.Multihash(buf)

	return mul.B58String(), nil
}
