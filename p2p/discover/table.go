package discover

import (
	mrand "math/rand"
	"net"
	"sync"
	"zeus/p2p/utils"

	"github.com/ethereum/go-ethereum/common"
)

const (
	hashBits = len(common.Hash{}) * 8
	nBuckets = hashBits / 15 // Number of buckets
)

type Table struct {
	mutex   sync.Mutex
	buckets [nBuckets]*bucket
	nursery []*Node
	rand    *mrand.Rand
	ips     utils.DistinctNetSet

	db         *nodeDB
	refreshReq chan chan struct{}
	initDone   chan struct{}
	closeReq   chan struct{}
	closed     chan struct{}

	bondmu    sync.Mutex
	bonding   map[NodeID]*bondproc
	bondslots chan struct{}

	net  transport
	self *Node
}

type bucket struct {
	entries      []*Node
	replacements []*Node
	ips          utils.DistinctNetSet
}

type bondproc struct {
	err  error
	n    *Node
	done chan struct{}
}

// transport is implemented by the UDP transport.
// it is an interface so we can test without opening lots of UDP
// sockets and without generating a private key.
type transport interface {
	ping(NodeID, *net.UDPAddr) error
	waitping(NodeID) error
	findnode(toid NodeID, addr *net.UDPAddr, target NodeID) ([]*Node, error)
	close()
}

func newTable(t transport, ourID NodeID, ourAddr *net.UDPAddr, nodeDBPath string, bootnodes []*Node) (*Table, error) {
	db, err := newNodeDB(nodeDBPath, ver)
}
