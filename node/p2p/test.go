package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

const messageID = 0

type Message string

var p = make(chan *p2p.Peer)

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	for {
		fmt.Println("11")
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage [1]Message
		err = msg.Decode(&myMessage)
		if err != nil {
			continue
		}

		switch myMessage[0] {
		case "foo":
			err := p2p.SendItems(ws, messageID, "bar")
			if err != nil {
				return err
			}
		default:
			fmt.Println("rece:", myMessage)
		}

	}
}

func MyProtocal() p2p.Protocol {
	return p2p.Protocol{
		Name:    "YourNode",
		Version: 1,
		Length:  1,
		Run:     msgHandler,
	}
}

func main() {
	nodeKey, _ := crypto.GenerateKey()
	config := p2p.Config{
		MaxPeers:   10,
		PrivateKey: nodeKey,
		Name:       "Test Node",
		ListenAddr: ":30300",
		Protocols:  []p2p.Protocol{MyProtocal()},
	}
	srv := p2p.Server{
		Config: config,
	}

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer srv.Stop()
	// fmt.Println("started...", srv.NodeInfo())

	coon, err := net.DialTimeout("tcp", srv.ListenAddr, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer coon.Close()

	select {
	case peer := <-p:
		if peer.LocalAddr().String() != coon.LocalAddr().String() {
			fmt.Println("wrong")
		}
		peers := srv.Peers()
		if !reflect.DeepEqual(peers, []*p2p.Peer{peer}) {
			fmt.Println("wrong2")
		}
		fmt.Println("YES")
	case <-time.After(5 * time.Second):
		fmt.Println("over time")
	}
}
