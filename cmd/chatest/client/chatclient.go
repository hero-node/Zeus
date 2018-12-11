package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"

	//"net/http"
	"net/url"
	"errors"
	crand "crypto/rand"
)

type Friend struct {
	address common.Address
	symID string
	pub *ecdsa.PublicKey
}

func newFriend(address string) *Friend {
	f :=  &Friend{address:common.HexToAddress(address), symID:""}
	return f
}

type SubscribeParam struct {
	Subscrib bool `json:"subscribe"`
	From string `json:"from"`
	To string `json:"to"`
}

type ReceiveMessage struct {
	Payload []byte `json:"payload"`
	From string `json:"from"`
	Encrypted bool `json:"encrypted"`
	Pub []byte `json:"pub"`
}

type NewMessage struct {
	From string `json:"from"`
	To string `json:"to"`
	Payload []byte `json:"payload"`   // encrypted or unencrypted message
	Encrypted bool `json:"encrypted"`
	Pub []byte `json:"pub""`
}

var privateKey *ecdsa.PrivateKey
var address common.Address

var friends [2]*Friend

var mytopic = []byte{0x11, 0xff, 0xdd, 0xaa}
var addr string
const minPow = 0.5
const powTime = 3

func main() {
	clientType := flag.Int("client", 0, "1 or 2 or 3")
	flag.Parse()
	fmt.Println("输入想要连接的Hero Node节点, 默认127.0.0.1:9000")

	fmt.Scanf("%s", &addr)
	if addr == "" {
		addr = "127.0.0.1:9000"
	}
	if *clientType == 0 {
		// 0x164732Dc9261b06B2C3bc700f1C534C999088585
		privateKey, _ = crypto.HexToECDSA("DE7559B1B70BF4F7567A24539DD79F4B072E3D7E10C7A4165AE9CE0F644E5F49")
		friends = [2]*Friend{newFriend("0x3a370712b70Ed656E762F218b1fa26e28CEcDd4d"), newFriend("0x6341Cb1797948be3352778057a5F616BC70feB62"),}
	} else if *clientType == 1 {
		// 0x3a370712b70Ed656E762F218b1fa26e28CEcDd4d
		privateKey, _ = crypto.HexToECDSA("76E47062AD13BE8759EFFA773A070DBC009424194EAC77C73DFB523EB5EB4C8A")
		friends = [2]*Friend{newFriend("0x6341Cb1797948be3352778057a5F616BC70feB62"), newFriend("0x164732Dc9261b06B2C3bc700f1C534C999088585"),}
	} else if *clientType == 2 {
		// 0x6341Cb1797948be3352778057a5F616BC70feB62
		privateKey, _ = crypto.HexToECDSA("C4AB8175F0BB7DDE95C91192E4123B913875DEAAC68DE9AEE309B55E8A1A2AE5")
		friends = [2]*Friend{newFriend("0x3a370712b70Ed656E762F218b1fa26e28CEcDd4d"), newFriend("0x164732Dc9261b06B2C3bc700f1C534C999088585"),}
	}
	address = crypto.PubkeyToAddress(privateKey.PublicKey)

	fmt.Println("此账户地址为：", address.String())


	wsUrl := url.URL{Scheme:"ws", Host:addr, Path:"subscrib"}
	c, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go func() {
		// write
		for _, f := range friends {
			sub := SubscribeParam{
				From: address.String(),
				To: f.address.String(),
				Subscrib: true,
			}

			jsonSend, _ := json.Marshal(sub)
			fmt.Println("准备建立和好友的连接... " + f.address.String())
			if err := c.WriteMessage(websocket.TextMessage, jsonSend); err == nil {
				fmt.Println("连接成功")
			}
		}
		for {
			echo()
		}

	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("挂了")
			c.Close()
			panic(err)
		}
		var mess ReceiveMessage
		if err = json.Unmarshal(message, &mess); err == nil {
			if mess.From != address.String() {
				for _, f := range friends {
					if f.address.String() == mess.From {
						f.pub, _ = crypto.UnmarshalPubkey(mess.Pub)
					}
				}
				if mess.Encrypted {
					fmt.Printf("解密前[%s]%s\n", mess.From[:6], string(mess.Payload))
					payload, _ := decrypt(mess.Payload, privateKey)
					fmt.Printf("解密后[%s]%s\n", mess.From[:6], string(payload))
				} else {
					fmt.Printf("未加密[%s]%s\n", mess.From[:6], string(mess.Payload))
				}
			}
		}
	}
}

func echo() {
	fmt.Println("选择好友发送")
	fmt.Println("0. " + friends[0].address.String())
	fmt.Println("1. " + friends[1].address.String())
	var i int
	fmt.Scanf("%d\n", &i)
	showSession(friends[i])
}

func showSession(f *Friend)  {
	fmt.Printf("已选择%s进行会话, 退出会话请输入quit并按回车\n", f.address.String())
	for {
		var s string
		fmt.Scanf("%s\n", &s)
		if s == "quit" {
			fmt.Println("\n\n")
			echo()
		}
		sendMessage(f, s)
	}
}

func sendMessage(f *Friend, message string) {
	httpUrl := url.URL{Scheme:"http", Host:addr, Path:"post"}
	fmt.Println("send: " + message)
	payload := []byte(message)
	if f.pub != nil {
		// encrypted
		ct, err := encrypt([]byte(message), f.pub)
		if err != nil {
			fmt.Println(err)
			return
		}
		payload = ct
	}

	m := NewMessage{
		From: address.String(),
		To: f.address.String(),
		Payload: payload,
		Encrypted: f.pub != nil,
		Pub: crypto.FromECDSAPub(&privateKey.PublicKey),
	}
	j, _ := json.Marshal(m)
	resp, err := http.Post(httpUrl.String(), "application/json", bytes.NewBuffer(j))
	if err != nil {
		fmt.Println("发送消息失败")
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
}

func encrypt(data []byte, pub *ecdsa.PublicKey) ([]byte, error)  {
	if !whisperv6.ValidatePublicKey(pub) {
		return []byte{}, errors.New("invalid public key provided for asymmetric encryption")
	}
	return ecies.Encrypt(crand.Reader, ecies.ImportECDSAPublic(pub), data, nil, nil)
}

func decrypt(ct []byte, key *ecdsa.PrivateKey) ([]byte, error)  {
	return ecies.ImportECDSA(key).Decrypt(ct, nil, nil)
}
