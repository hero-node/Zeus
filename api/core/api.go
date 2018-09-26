package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
	"zeus/api/noderror"
	"zeus/utils/config"
	"zeus/utils/global"

	"context"

	"errors"
	"strings"

	"math/big"
	"runtime"

	"zeus/api/bootstrap"
	"zeus/api/utils"
	"zeus/nameservice/ens"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var ethHost string

func InitRoute(router *gin.Engine) {
	ethHost = config.GetValidEthHost()

	router.GET("/isHero", isHero)
	router.GET("/available/:chain", getChainAvailable)
	router.GET("/balance/:chain/:address", getBalance)
	router.GET("/chains", getAllAvailableChains)
	router.GET("/gasPrice/:chain", getGasPrice)
	router.GET("/gasPrice", getEthGasPrice)
	router.GET("/gasLimit/:chain", getGasLimit)
	router.GET("/gasLimit", getEthGasLimit)
	router.GET("/info", getNodeInfo)
	// router.GET("/mining/:chain", getMining)
	router.GET("/filterLogs/:id", getEthFilterLogs)

	router.GET("/block/:chain", getBlockByHeightOrHash)
	router.GET("/blockHash/:chain/:height", getBlockHashByHeight)
	router.GET("/blockHeight/:chain", getBlockHeight)

	router.GET("/account/transaction/count/:chain/:hash", getTransactionCountInBlock)
	router.GET("/sendRawTransaction/:chain/:data", sendRawTransaction)
	router.GET("/transaction/:chain/:hash", getTransactionByHash)
	router.GET("/transactionReceipt/:chain/:hash", getReceiptByHash)
	router.GET("/peers", getPeers)

	// ens
	router.GET("/ens/:ensname", ensParse)
	router.GET("/ensEncode/:content", ensEncode)

	// ipfs
	router.POST("/ipfs/add", ReverseProxy())
	router.GET("/ipfs/bitswap/ledger", ReverseProxy())
	router.GET("/ipfs/bitswap/reprovide", ReverseProxy())
	router.GET("/ipfs/bitswap/stat", ReverseProxy())
	router.GET("/ipfs/bitswap/unwant", ReverseProxy())
	router.GET("/ipfs/bitswap/wantlist", ReverseProxy())
	router.GET("/ipfs/block/get", ReverseProxy())
	router.POST("/ipfs/block/put", ReverseProxy())
	router.GET("/ipfs/block/rm", ReverseProxy())
	router.GET("/ipfs/block/stat", ReverseProxy())
	router.GET("/ipfs/bootstrap/add/default", ReverseProxy())
	router.GET("/ipfs/bootstrap/list", ReverseProxy())
	router.GET("/ipfs/bootstrap/rm/all", ReverseProxy())
	router.GET("/ipfs/cat", ReverseProxy())
	router.GET("/ipfs/commands", ReverseProxy())
	router.GET("/ipfs/config/edit", ReverseProxy())
	router.GET("/ipfs/config/replace", ReverseProxy())
	router.GET("/ipfs/config/show", ReverseProxy())
	router.GET("/ipfs/dag/get", ReverseProxy())
	router.POST("/ipfs/dag/put", ReverseProxy())
	router.GET("/ipfs/dag/resolve", ReverseProxy())
	router.GET("/ipfs/dht/findpeer", ReverseProxy())
	router.GET("/ipfs/dht/findprovs", ReverseProxy())
	router.GET("/ipfs/dht/get", ReverseProxy())
	router.GET("/ipfs/dht/provide", ReverseProxy())
	router.POST("/ipfs/dht/put", ReverseProxy())
	router.GET("/ipfs/dht/query", ReverseProxy())
	router.GET("/ipfs/diag/cmds/clear", ReverseProxy())
	router.GET("/ipfs/diag/cmds/set-time", ReverseProxy())
	router.GET("/ipfs/diag/sys", ReverseProxy())
	router.GET("/ipfs/dns", ReverseProxy())
	router.GET("/ipfs/file/ls", ReverseProxy())
	router.GET("/ipfs/files/cp", ReverseProxy())
	router.GET("/ipfs/files/flush", ReverseProxy())
	router.GET("/ipfs/files/ls", ReverseProxy())
	router.GET("/ipfs/files/mkdir", ReverseProxy())
	router.GET("/ipfs/files/mv", ReverseProxy())
	router.GET("/ipfs/files/read", ReverseProxy())
	router.GET("/ipfs/files/rm", ReverseProxy())
	router.GET("/ipfs/files/stat", ReverseProxy())
	router.GET("/ipfs/files/write", ReverseProxy())
	router.GET("/ipfs/filestore/dups", ReverseProxy())
	router.GET("/ipfs/filestore/ls", ReverseProxy())
	router.GET("/ipfs/filestore/verify", ReverseProxy())
	router.GET("/ipfs/get", ReverseProxy())
	router.GET("/ipfs/id", ReverseProxy())
	router.GET("/ipfs/key/gen", ReverseProxy())
	router.GET("/ipfs/key/list", ReverseProxy())
	router.GET("/ipfs/key/rename", ReverseProxy())
	router.GET("/ipfs/key/rm", ReverseProxy())
	router.GET("/ipfs/log/level", ReverseProxy())
	router.GET("/ipfs/log/ls", ReverseProxy())
	router.GET("/ipfs/log/tail", ReverseProxy())
	router.GET("/ipfs/ls", ReverseProxy())
	router.GET("/ipfs/mount", ReverseProxy())
	router.GET("/ipfs/name/publish", ReverseProxy())
	router.GET("/ipfs/name/resolve", ReverseProxy())
	router.GET("/ipfs/object/data", ReverseProxy())
	router.GET("/ipfs/object/diff", ReverseProxy())
	router.GET("/ipfs/object/get", ReverseProxy())
	router.GET("/ipfs/object/links", ReverseProxy())
	router.GET("/ipfs/object/new", ReverseProxy())
	router.GET("/ipfs/object/patch/add-link", ReverseProxy())
	router.GET("/ipfs/object/patch/append-data", ReverseProxy())
	router.GET("/ipfs/object/patch/rm-link", ReverseProxy())
	router.GET("/ipfs/object/patch/set-data", ReverseProxy())
	router.POST("/ipfs/object/put", ReverseProxy())
	router.GET("/ipfs/object/stat", ReverseProxy())
	router.GET("/ipfs/p2p/listener/close", ReverseProxy())
	router.GET("/ipfs/p2p/listener/ls", ReverseProxy())
	router.GET("/ipfs/p2p/listener/open", ReverseProxy())
	router.GET("/ipfs/p2p/stream/close", ReverseProxy())
	router.GET("/ipfs/p2p/stream/dial", ReverseProxy())
	router.GET("/ipfs/p2p/stream/ls", ReverseProxy())
	router.GET("/ipfs/pin/add", ReverseProxy())
	router.GET("/ipfs/pin/ls", ReverseProxy())
	router.GET("/ipfs/pin/rm", ReverseProxy())
	router.GET("/ipfs/pin/update", ReverseProxy())
	router.GET("/ipfs/pin/verify", ReverseProxy())
	router.GET("/ipfs/ping", ReverseProxy())
	router.GET("/ipfs/pubsub/ls", ReverseProxy())
	router.GET("/ipfs/pubsub/peers", ReverseProxy())
	router.GET("/ipfs/pubsub/pub", ReverseProxy())
	router.GET("/ipfs/pubsub/sub", ReverseProxy())
	router.GET("/ipfs/refs/local", ReverseProxy())
	router.GET("/ipfs/repo/fsck", ReverseProxy())
	router.GET("/ipfs/repo/gc", ReverseProxy())
	router.GET("/ipfs/repo/stat", ReverseProxy())
	router.GET("/ipfs/repo/verify", ReverseProxy())
	router.GET("/ipfs/repo/version", ReverseProxy())
	router.GET("/ipfs/resolve", ReverseProxy())
	router.GET("/ipfs/shutdown", ReverseProxy())
	router.GET("/ipfs/stats/bitswap", ReverseProxy())
	router.GET("/ipfs/stats/bw", ReverseProxy())
	router.GET("/ipfs/stats/repo", ReverseProxy())
	router.GET("/ipfs/swarm/addrs/listen", ReverseProxy())
	router.GET("/ipfs/swarm/addrs/local", ReverseProxy())
	router.GET("/ipfs/swarm/connect", ReverseProxy())
	router.GET("/ipfs/swarm/disconnect", ReverseProxy())
	router.GET("/ipfs/swarm/filters/add", ReverseProxy())
	router.GET("/ipfs/swarm/filters/rm", ReverseProxy())
	router.GET("/ipfs/swarm/peers", ReverseProxy())
	router.POST("/ipfs/tar/add", ReverseProxy())
	router.GET("/ipfs/tar/cat", ReverseProxy())
	router.GET("/ipfs/update", ReverseProxy())
	router.GET("/ipfs/version", ReverseProxy())

}

func isHero(c *gin.Context) {
	c.JSON(200, gin.H{
		"result":  "success",
		"content": 1,
	})
}

func getChainAvailable(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))

	ok := false
	switch chain {
	case global.BTC:
		// TODO: BTC
	case global.ETH:
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := ethclient.Dial(ethHost)
		if err != nil {
			ok = false
		} else {
			ok = true
		}
	case global.QTUM:
		_, err := Call_QTUM("getinfo", []interface{}{})
		if err != nil {
			ok = false
		} else {
			ok = true
		}

	default:
	}

	if ok {
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"available": 1,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"available": 0,
			},
		})
	}
}

func getBalance(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	address := c.Param("address")

	switch chain {
	case global.BTC:
		// TODO:
		c.JSON(500, gin.H{
			"result": "error",
		})
	case global.ETH:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		balance, err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"chain":   "eth",
				"balance": balance,
			},
		})

	case global.QTUM:
		// TODO:
		c.JSON(500, gin.H{
			"result": "error",
		})
	default:
		c.JSON(500, gin.H{
			"result": "error",
		})
	}
}

func getAllAvailableChains(c *gin.Context) {
	rsp := gin.H{}
	// TODO: BTC

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.Dial(ethHost)
	if err == nil {
		networkID, err := client.NetworkID(ctx)
		if err == nil {
			rsp["eth"] = gin.H{
				"networkID": networkID,
			}
		}
	}

	qtumresp, err := Call_QTUM("getnetworkinfo", []interface{}{})
	if err == nil {
		rsp["qtum"] = gin.H{
			"networkVersion": qtumresp.Result.(map[string]interface{})["version"],
		}
	}

	c.JSON(200, rsp)
}

func getGasLimit(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))

	isEth := false

	from := c.Query("from")
	to := c.Query("to")
	data := c.Query("data")

	obj := map[string]string{}
	if len(from) > 0 {
		obj["from"] = from
	}
	if len(to) > 0 {
		obj["to"] = to
	}
	if len(data) > 0 {
		obj["data"] = data
	}

	switch chain {
	case global.BTC:
		c.JSON(500, gin.H{
			"result": "error",
			"reason": "No gas for btc",
		})
	case global.ETH:
		isEth = true

	case global.QTUM:
		c.JSON(500, gin.H{
			"result": "error",
			"reason": "No gas for qtum",
		})
	default:
		isEth = true
	}

	if isEth {
		resp, err := Call_ETH("eth_estimateGas", []interface{}{obj})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"chain":    "eth",
				"gasLimit": resp.Result.(map[string]interface{})["result"],
			},
		})
	}
}

func getEthGasLimit(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")
	data := c.Query("data")

	obj := map[string]string{}
	if len(from) > 0 {
		obj["from"] = from
	}
	if len(to) > 0 {
		obj["to"] = to
	}
	if len(data) > 0 {
		obj["data"] = data
	}

	resp, err := Call_ETH("eth_estimateGas", []interface{}{obj})
	if err != nil {
		noderror.Error(err, c)
		return
	}

	c.JSON(200, gin.H{
		"result": "success",
		"content": gin.H{
			"chain":    "eth",
			"gasLimit": resp.Result.(map[string]interface{})["result"],
		},
	})
}

func getGasPrice(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))

	isEth := false

	switch chain {
	case global.BTC:
		c.JSON(500, gin.H{
			"result": "error",
			"reason": "No gas for btc",
		})
	case global.ETH:
		isEth = true

	case global.QTUM:
		c.JSON(500, gin.H{
			"result": "error",
			"reason": "No gas for qtum",
		})
	default:
		isEth = true
	}

	if isEth {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"chain":    "eth",
				"gasPrice": gasPrice,
			},
		})
	}
}

func getEthGasPrice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := ethclient.Dial(ethHost)
	if err != nil {
		noderror.Error(err, c)
		return
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		noderror.Error(err, c)
		return
	}
	c.JSON(200, gin.H{
		"result": "success",
		"content": gin.H{
			"chain":    "eth",
			"gasPrice": gasPrice,
		},
	})
}

func getNodeInfo(c *gin.Context) {
	info := "HN-v0.0.1-debug/" + runtime.GOOS + "-" + runtime.GOARCH + "/" + runtime.Version()
	c.JSON(200, gin.H{
		"result":  "success",
		"content": info,
	})
}

func getEthFilterLogs(c *gin.Context) {
	id := c.Param("id")

	resp, err := Call_ETH("eth_getFilterLogs", []interface{}{id})
	if err != nil {
		noderror.Error(err, c)
		return
	}

	c.JSON(200, gin.H{
		"result":  "success",
		"content": resp.Result,
	})
}

// func getMining(c *gin.Context) {
// 	chain := strings.ToLower(c.Param("chain"))

// 	switch chain {
// 	case global.BTC:
// 		// TODO
// 	case global.ETH:
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		client, err := ethclient.Dial(ethHost)
// 		if err != nil {
// 			noderror.Error(err, c)
// 			return
// 		}
// 		mining, err := client.
// 	case global.QTUM:
// 		// TODO
// 	default:
// 	}
// }

func getBlockByHeightOrHash(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	height := c.Query("height")
	hash := c.Query("hash")

	switch chain {
	case global.BTC:
		// TODO
	case global.ETH:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		if len(height) > 0 {
			heightInt := new(big.Int)
			heightInt, ok := heightInt.SetString(height, 10)
			if !ok {
				noderror.Error(errors.New("Parameters"), c)
				return
			}
			block, err := client.BlockByNumber(ctx, heightInt)
			if err != nil {
				noderror.Error(err, c)
				return
			}
			c.JSON(200, gin.H{
				"result":  "success",
				"content": block.Header(),
			})
		} else if len(hash) > 0 {
			block, err := client.BlockByHash(ctx, common.HexToHash(hash))
			if err != nil {
				noderror.Error(err, c)
				return
			}
			c.JSON(200, gin.H{
				"result":  "success",
				"content": block.Header(),
			})
		}
	case global.QTUM:
		if len(height) > 0 {
			heightInt, err := strconv.Atoi(height)
			if err != nil {
				noderror.Error(err, c)
				return
			}

			resp, err := Call_QTUM("getblockhash", []interface{}{heightInt})
			if err != nil {
				noderror.Error(err, c)
				return
			}

			hash = resp.Result.(string)
		}

		resp, err := Call_QTUM("getblock", []interface{}{hash})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		block := resp.Result
		c.JSON(200, gin.H{
			"result":  "success",
			"content": block,
		})
	}
}

func getBlockHashByHeight(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	height := c.Param("height")

	switch chain {
	case global.BTC:
		// TODO
	case global.ETH:
		heightInt := new(big.Int)
		heightInt, ok := heightInt.SetString(height, 10)
		if !ok {
			noderror.Error(errors.New("Parameters Error"), c)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		block, err := client.BlockByNumber(ctx, heightInt)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"chain": "eth",
				"hash":  block.Hash(),
			},
		})

	case global.QTUM:
		heightInt, err := strconv.Atoi(height)
		if err != nil {
			noderror.Error(err, c)
			return
		}

		resp, err := Call_QTUM("getblockhash", []interface{}{heightInt})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		hash := resp.Result.(string)
		c.JSON(200, gin.H{
			"result":  "success",
			"content": hash,
		})
	}
}

func getBlockHeight(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))

	switch chain {
	case global.BTC:
		// TODO
	case global.ETH:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		block, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"chain":  "eth",
				"height": block.Number(),
			},
		})

	case global.QTUM:
		resp, err := Call_QTUM("getblockcount", []interface{}{})
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result":  "success",
			"content": resp.Result,
		})
	}
}

func getTransactionCountInBlock(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	hash := c.Param("hash")

	switch chain {
	case global.BTC:
		// TODO
	case global.ETH:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := ethclient.Dial(ethHost)
		if err != nil {
			noderror.Error(err, c)
			return
		}
		count, err := client.TransactionCount(ctx, common.HexToHash(hash))
		if err != nil {
			noderror.Error(err, c)
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": gin.H{
				"count": count,
			},
		})
	case global.QTUM:
		// TODO
	}
}

func sendRawTransaction(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	data := c.Param("data")

	switch chain {
	case global.BTC:
		// TODO
	case global.ETH:
		resp, err := Call_ETH("eth_sendRawTransaction", []interface{}{data})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		c.JSON(200, gin.H{
			"result":  "success",
			"content": resp,
		})

	case global.QTUM:

	}
}

func getTransactionByHash(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	hash := c.Param("hash")

	switch chain {
	case global.BTC:

	case global.ETH:
		resp, err := Call_ETH("eth_getTransactionByHash", []interface{}{hash})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		c.JSON(200, gin.H{
			"result":  "success",
			"content": resp,
		})

	case global.QTUM:

	}
}

func getReceiptByHash(c *gin.Context) {
	chain := strings.ToLower(c.Param("chain"))
	hash := c.Param("hash")

	switch chain {
	case global.BTC:

	case global.ETH:
		resp, err := Call_ETH("eth_getTransactionReceipt", []interface{}{hash})
		if err != nil {
			noderror.Error(err, c)
			return
		}

		c.JSON(200, gin.H{
			"result":  "success",
			"content": resp,
		})

	case global.QTUM:

	}
}

var max int

func getPeers(c *gin.Context) {
	localPath := "http://localhost" + config.GetHttpPort() + "/ipfs/swarm/peers"
	resp, err := http.Get(localPath)
	if err != nil {
		noderror.Error(err, c)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		noderror.Error(err, c)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		noderror.Error(err, c)
		return
	}

	peers := result["Peers"].([]interface{})
	addrs := make([]string, 0)
	for _, peer := range peers {
		addr := utils.ExtractIPv4(peer.(map[string]interface{})["Addr"].(string))
		dup := false
		if addr != "" {
			for _, a := range addrs {
				if a == addr {
					dup = true
				}
			}

			if !dup {
				addrs = append(addrs, addr)
			}
		}
	}

	addrs = append(addrs, bootstrap.B.Bootlist...)

	max = len(addrs)
	mutex := sync.Mutex{}
	done := make(chan int)
	response := []string{}

	for _, a := range addrs {

		go func(addr string) {
			path := utils.ConstructUrl(addr)
			netClient := http.Client{Timeout: time.Second * 5}
			resp, err := netClient.Get(path + "/isHero")

			if err != nil {
				mutex.Lock()
				max = max - 1
				if len(response) == max {
					done <- 1
				}
				mutex.Unlock()
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				mutex.Lock()
				max = max - 1
				if len(response) == max {
					done <- 1
				}
				mutex.Unlock()
				return
			}
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				mutex.Lock()
				max = max - 1
				if len(response) == max {
					done <- 1
				}
				mutex.Unlock()
				return
			}
			if result["result"] == "success" {
				mutex.Lock()
				response = append(response, addr)

				if len(response) == max {
					done <- 1
				}
				mutex.Unlock()
			} else {
				mutex.Lock()
				max = max - 1

				if len(response) == max {
					done <- 1
				}
				mutex.Unlock()
			}
		}(a)
	}

	select {
	case <-done:
		c.JSON(200, gin.H{
			"result":  "success",
			"content": response,
		})

	case <-time.Tick(time.Second * 12):
		defer func() { <-done }()
		c.JSON(500, gin.H{
			"result": "error",
			"reason": "Timeout",
		})
	}
}

func ensParse(c *gin.Context) {
	// 1.find mapping content
	// 2.decode to ipfs hash
	// 3.redirect to ipfs resource url
	ensName := c.Param("ensname")
	content := ens.EnsToContent(ensName)
	ipfsHash, err := ens.IpfsDecode(content)
	if err != nil {
		noderror.Error(err, c)
		return
	}

	dest := "https://ipfs.io/ipfs/" + ipfsHash

	c.Redirect(http.StatusMovedPermanently, dest)
}

func ensEncode(c *gin.Context) {
	ipfsHash := c.Param("content")
	ensHash, err := ens.IpfsEncode(ipfsHash)
	if err != nil {
		noderror.Error(err, c)
		return
	}
	c.JSON(200, gin.H{
		"result":  "success",
		"content": ensHash,
	})
}
