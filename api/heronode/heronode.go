package heronode

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"zeus/api/noderror"
	"zeus/utils/global"

	"context"

	"errors"
	"strings"

	"math/big"
	"runtime"

	"zeus/api/bootstrap"
	"zeus/api/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var ethHost string

func InitRoute(router *gin.Engine) {
	ethHost = global.Ethhost()

	router.GET("/isHero", isHero)
	router.GET("/available/:chain", getChainAvailable)
	router.GET("/balance/:chain/:address", getBalance)
	router.GET("/chains", getAllAvailableChains)
	router.GET("/gasPrice/:chain", getGasPrice)
	router.GET("/gasPrice", getEthGasPrice)
	router.GET("/info", getNodeInfo)
	// router.GET("/mining/:chain", getMining)

	router.GET("/block/:chain", getBlockByHeightOrHash)
	router.GET("/blockHash/:chain/:height", getBlockHashByHeight)
	router.GET("/blockHeight/:chain", getBlockHeight)

	router.GET("/account/transaction/count/:chain/:hash", getTransactionCountInBlock)
	router.GET("/sendRawTransaction/:chain/:data", sendRawTransaction)
	router.GET("/transaction/:chain/:hash", getTransactionByHash)
	router.GET("/transactionReceipt/:chain/:hash", getReceiptByHash)
	router.GET("/peers", getPeers)

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
		// TODO: QTUM
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

	// TODO: QTUM

	c.JSON(200, rsp)
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
		// TODO
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
		// TODO
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

func getPeers(c *gin.Context) {
	// TODO: no hard code
	resp, err := http.Get("http://localhost:8080/ipfs/swarm/peers")
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

	addrs = append(addrs, bootstrap.B.bootlist)

	max := len(addrs)
	mutex := sync.Mutex{}
	done := make(chan int)
	response := []string{}
	for _, a := range addrs {
		path := utils.ConstructUrl(a)
		go func() {
			resp, err := http.Get(path + "/isHero")
			if err != nil {
				mutex.Lock()
				max = max - 1
				mutex.Unlock()
				if len(response) == max {
					done <- 1
				}
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				mutex.Lock()
				max = max - 1
				mutex.Unlock()
				if len(response) == max {
					done <- 1
				}
				return
			}
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				mutex.Lock()
				max = max - 1
				mutex.Unlock()
				if len(response) == max {
					done <- 1
				}
				return
			}
			if result["result"] == "success" {
				mutex.Lock()
				response = append(response, a)
				mutex.Unlock()
				if len(response) == max {
					done <- 1
				}
			} else {
				mutex.Lock()
				max = max - 1
				mutex.Unlock()
				if len(response) == max {
					done <- 1
				}
			}
		}()
	}

	<-done
	c.JSON(200, gin.H{
		"result":  "success",
		"content": response,
	})
}
