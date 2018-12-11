package main

import (
	"context"
	"crypto/ecdsa"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

var shh *shhclient.Client

type NewMessage struct {
	From string `json:"from"`
	To string `json:"to"`
	Payload []byte `json:"payload"`   // encrypted or unencrypted message
	Encrypted bool `json:"encrypted"`
	Pub []byte `json:"pub""`
}


func (message NewMessage)ToWhisperType() (whisperv6.NewMessage, error)  {
	var wMessage = whisperv6.NewMessage{}
	pass := GeneratePassword(message.From, message.To)
	symID, err := GetSymKeyFromPassword(pass)
	if err != nil {
		return wMessage, err
	}

	var flagByte byte = 0x0
	if message.Encrypted {
		flagByte = 0x1
	}
	fromBytes := []byte(message.From)

	payload := make([]byte, 0)
	payload = append(payload, flagByte)
	payload = append(payload, fromBytes...)
	payload = append(payload, message.Pub...)
	payload = append(payload, message.Payload...)

	wMessage = whisperv6.NewMessage{
		SymKeyID: symID,
		Topic: whisperv6.BytesToTopic([]byte{0x11, 0x22, 0x33, 0x44}),
		Payload: payload,
		PowTime: 3,
		PowTarget: 0.5,
	}
	return wMessage, nil
}

type ReceiveMessage struct {
	Payload []byte `json:"payload"`
	From string `json:"from"`
	Encrypted bool `json:"encrypted"`
	Pub []byte `json:"pub"`
}

func NewReceiveMessageFromWhisper(message whisperv6.Message) (ReceiveMessage)  {
	payload := message.Payload
	fmt.Println("payload len ", len(message.Payload))
	encryptedBytes := payload[0]
	fromBytes := payload[1:43]
	pubBytes := payload[43:108]
	contentBytes := payload[108:]
	var b = false
	if encryptedBytes == 0x1 {
		b = true
	}

	fmt.Println("Rece From: ", string(fromBytes))
	return ReceiveMessage {
		Payload: contentBytes,
		From: string(fromBytes),
		Encrypted: b,
		Pub: pubBytes,
	}

}

type SubscribeParam struct {
	Subscrib bool `json:"subscribe"`
	From string `json:"from"`
	To string `json:"to"`
}

func (s SubscribeParam)ToWhisperType() (whisperv6.Criteria)  {
	pass := GeneratePassword(s.From, s.To)
	symID, _ := GetSymKeyFromPassword(pass)
	return whisperv6.Criteria{
		SymKeyID: symID,
		MinPow: 0.5,
		Topics:[]whisperv6.TopicType{whisperv6.BytesToTopic([]byte{0x11, 0x22, 0x33, 0x44})},
	}
}

func GeneratePassword(add1 string, add2 string) (string) {
	if strings.Compare(add1, add2) >= 0 {
		return add1 + add2;
	}
	return add2 + add1;
}

func GetSymKeyFromPassword(password string) (string, error)  {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return shh.GenerateSymmetricKeyFromPassword(ctx, password)
}

func main()  {
	var upgrader = websocket.Upgrader {
		ReadBufferSize:1024,
		WriteBufferSize:1024,
	}

	port := flag.String("port", ":9000", "port")
	flag.Parse()


	fmt.Println("输入想连接的whisper节点url, 默认ws://127.0.0.1:8021")
	var whisperUrl string
	fmt.Scanf("%s\n", &whisperUrl)
	if whisperUrl == "" {
		whisperUrl = "ws://127.0.0.1:8021"
	}
	shhh, err := shhclient.Dial(whisperUrl)

	if err != nil {
		panic(err)
	}
	shh = shhh

	r := gin.Default()
	r.GET("/symKeyFromPassword/:password", func(c *gin.Context) {
		password := c.Param("password")
		symID, err := GetSymKeyFromPassword(password)
		if err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"reason": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"result": "success",
			"content": symID,
		})
	})

	r.POST("/post", func(c *gin.Context) { // content-type must application/json
		var message NewMessage

		if err := c.BindJSON(&message); err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"reason": err.Error(),
			})
			return
		}


		wMessage, err := message.ToWhisperType()
		if err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"reason": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		hash, err := shh.Post(ctx, wMessage)
		if err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"reason": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"result": "success",
			"content": hash,
		})
	})

	r.GET("/subscrib", func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"reason": err.Error(),
			})
			return
		}


		mess := make(chan *whisperv6.Message)

		go func() {
			// read
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					conn.SetWriteDeadline(time.Now().Add(time.Second*10))
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					conn.Close()
					return
				}
				var subscribeParam SubscribeParam
				if err = json.Unmarshal(message, &subscribeParam); err == nil {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					crit := subscribeParam.ToWhisperType()
					shh.SubscribeMessages(ctx, crit, mess)
					fmt.Println("订阅 ", crit.SymKeyID)
				}
			}

		}()

		go func() {
			// write
			for {
				select {
				case message, ok := <- mess:
					fmt.Println("send mess")
					conn.SetWriteDeadline(time.Now().Add(time.Second*10))
					if ok {
						rMessage := NewReceiveMessageFromWhisper(*message)
						messageBytes, _ := json.Marshal(rMessage)
						err := conn.WriteMessage(websocket.TextMessage, messageBytes)
						if err != nil {
							conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
							conn.Close()
						}
					}
				}
			}
		}()
	})

	r.Run(*port)

}

func encrypt(data []byte, pub *ecdsa.PublicKey) ([]byte, error)  {
	if !whisperv6.ValidatePublicKey(pub) {
		return []byte{}, errors.New("invalid public key provided for asymmetric encryption")
	}
	return ecies.Encrypt(crand.Reader, ecies.ImportECDSAPublic(pub), data, nil, nil)
}