package hernode

import (
	l "log"
	"net"
	"net/rpc/jsonrpc"
	"time"
	"zeus/node/noderror"
	"zeus/node/rpcobjc"
	"zeus/utils/global"

	"context"

	"errors"
	"strings"

	"math/big"
	"runtime"

	"github.com/Mercy-Li/Goconfig/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var ethHost string

func InitRoute(router *gin.Engine) {
	ethHostLocal, err := config.GetConfigString("ethnode")
	ethHost = ethHostLocal
	if err != nil {
		l.Fatalln("get ethnode error")
	}

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
		client, err := net.DialTimeout("tcp", "106.14.187.240:8545", 10*time.Second)
		if err != nil {
			noderror.Error(err, c)
			return
		}

		clientRPC := jsonrpc.NewClient(client)
		var resp rpcobjc.ETHResp
		err = clientRPC.Call("eth_sendRawTransaction", data, &resp)
		if err != nil {
			noderror.Error(err, c)
			return
		}

		c.JSON(200, gin.H{
			"result":  "success",
			"content": resp,
		})

	case global.QTUM:
		// TODO:
	}
}
