package discover

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"net"
	"time"

	"zeus/p2p/nat"
	"zeus/p2p/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

const Version = 4

// Errors
var (
	errPacketTooSmall   = errors.New("too small")
	errBadHash          = errors.New("bad hash")
	errExpired          = errors.New("expired")
	errUnsolicitedReply = errors.New("unsolicited reply")
	errUnknownNode      = errors.New("unknown node")
	errTimeout          = errors.New("RPC timeout")
	errClockWarp        = errors.New("reply deadline too far in the future")
	errClosed           = errors.New("socket closed")
)

// Timeouts
const (
	respTimeout = 500 * time.Millisecond
	sendTimeout = 500 * time.Millisecond
	expiration  = 20 * time.Second

	ntpFailureThreshold = 32               // Continuous timeouts after which to check NTP
	ntpWarningCooldown  = 10 * time.Minute // Minimum amount of time to pass before repeating NTP warning
	driftThreshold      = 10 * time.Second // Allowed clock drift before warning user
)

// RPC packet types
const (
	pingPacket = iota + 1 // zero is 'reserved'
	pongPacket
	findnodePacket
	neighborsPacket
)

type (
	ping struct {
		Version    uint
		From, to   rpcEndpoint
		Expiration uint64

		Rest []rlp.RawValue `rlp:"tail"`
	}

	pong struct {
		To         rpcEndpoint
		ReplyTok   []byte
		Expiration uint64

		Rest []rlp.RawValue `rlp:"tail"`
	}

	findnode struct {
		Traget     NodeID
		Expiration uint64

		Rest []rlp.RawValue `rlp:"tail"`
	}

	neighbors struct {
		Nodes      []rpcNode
		Expiration uint64

		Rest []rlp.RawValue `rlp:"tail"`
	}

	rpcNode struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
		ID  NodeID
	}

	rpcEndpoint struct {
		IP  net.IP // len 4 for IPv4 or 16 for IPv6
		UDP uint16 // for discovery protocol
		TCP uint16 // for RLPx protocol
	}
)

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

type udp struct {
	conn        conn
	netrestrict utils.Netlist
	priv        *ecdsa.PrivateKey
	ourEndpoint rpcEndpoint

	addpending chan *pending
	gotreply   chan reply

	closing chan struct{}
	nat     nat.Interface

	*Table
}

type pending struct {
	from  NodeID
	ptype byte

	deadline time.Time

	callback func(resp interface{}) (done bool)
	errc     chan<- error
}

type reply struct {
	from  NodeID
	ptype byte

	data   interface{}
	mached chan<- bool
}

type ReadPacket struct {
	Data []byte
	Addr *net.UDPAddr
}

type Config struct {
	PrivateKey *ecdsa.PrivateKey

	/* optional */
	AnnounceAddr *net.UDPAddr
	NodedbPath   string
	NetRestrict  *utils.Netlist
	Bootnodes    []*Node
	Unhandled    chan<- ReadPacket
}

func ListenUDP(c conn, cfg Config) (*Table, error) {
	//	tab, _
}

func newUDP(c conn, cfg Config) (*Table, *udp, error) {
	udp := &udp{
		conn:        c,
		priv:        cfg.PrivateKey,
		netrestrict: cfg.NetRestrict,
		closing:     make(chan struct{}),
		gotreply:    make(chan reply),
		addpending:  make(chan *pending),
	}
	realaddress := c.LocalAddr().(*net.UDPAddr)
	udp.ourEndpoint = makeEndpoint(realaddress, uint16(realaddress.Port))
	tab, err := newTable(udp, PubkeyID(&cfg.PrivateKey.PublicKey), realaddress, cfg.NodedbPath, cfg.Bootnodes)
	if err != nil {
		return nil, nil, err
	}

	udp.Table = tab

	//	TODO: go udp.loop
}

func (t *udp) close() {
	close(t.closing)
	t.conn.Close()
}

func (t *udp) ping(toid NodeID, toaddr *net.UDPAddr) error {
	req := &ping{
		Version:    Version,
		From:       t.ourEndpoint,
		To:         makeEndpoint(toaddr, 0),
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	}
	packet, hash, err := encodePacket(t.priv, pingPacket, req)
	if err != nil {
		return err
	}

	errc := t.pending(toid, pongPacket, func(p interface{}) bool {
		return bytes.Equal(p.(*pong).ReplyTok, hash)
	})
	t.write(toaddr, req.name(), packet)
	return <-errc

}

func (t *udp) waitping(from NodeID) error {
	return <-t.pending(from, pingPacket, func(interface{}) bool { return true })
}

func (t *udp) pending(id NodeID, ptype byte, callback func(interface{}) bool) <-chan error {
	ch := make(chan error, 1)
	p := &pending{from: id, ptype: ptype, callback: callback, errc: ch}
	select {
	case t.addpending <- p:
	// loop will handle it

	case <-t.closing:
		ch <- errClosed
	}
	return ch
}

func (t *udp) write(toaddr *net.UDPAddr, what string, packet []byte) error {
	_, err := t.conn.WriteToUDP(packet, toaddr)
	log.Trace(">> "+what, "addr", toaddr, "err", err)
	return err
}

func makeEndpoint(addr *net.UDPAddr, tcpPort uint16) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

func (t *udp) nodeFromRPC(sender *net.UDPAddr, rn rpcNode) (*Node, error) {
	if rn.UDP <= 1024 {
		return nil, errors.New("low port")
	}
	if err := nat.CheckRelayIP(sender.IP, rn.IP); err != nil {
		return nil, err
	}
}

const (
	macSize  = 256 / 8
	sigSize  = 520 / 8
	headSize = macSize + sigSize // space of packet frame data
)

var (
	headSpace = make([]byte, headSize)

	// Neighbors replies are sent across multiple packets to
	// stay below the 1280 byte limit. We compute the maximum number
	// of entries by stuffing a packet until it grows too large.
	maxNeighbors int
)

func encodePacket(priv *ecdsa.PrivateKey, ptype byte, req interface{}) (packet, hash []byte, err error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)

	if err := rlp.Encode(b, req); err != nil {
		log.Error("Cannot encode packet")
		return nil, nil, err
	}
	packet := b.Bytes()
	sig, err := crypto.Sign(crypto.Keccak256(packet[headSize:]), priv)
	if err != nil {
		log.Error("Cannot encode packet")
		return nil, nil, err
	}

	copy(packet[macSize:], sig)
	hash := crypto.Keccak256(packet[macSize:])
	copy(packet, hash)
	return packet, hash, nil
}

func (req *ping) name() string { return "PING/v4" }
