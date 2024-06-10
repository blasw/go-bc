package network

import (
	"bytes"
	"fmt"
	"sync"
)

// LocalTransport implements Transport interface and is being used as a temporary local solution
type LocalTransport struct {
	addr      NetAddr
	peers     map[NetAddr]*LocalTransport
	lock      sync.RWMutex
	consumeCh chan RPC
}

// NewLocalTransport returns a new instance of LocalTransport
func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		peers:     make(map[NetAddr]*LocalTransport),
		consumeCh: make(chan RPC, 1024),
	}
}

// Consume returns a channel that can be used to consume RPC messages
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// Connect adds a peer to the list of connected peers
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

// SendMessage sends a message to a peer by address
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}

	return nil
}

// Broadcast sends a message to all connected peers
func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.addr, payload); err != nil {
			return err
		}
	}

	return nil
}

// Returns the address of the local transport
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
