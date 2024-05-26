package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	_ = tra.Connect(trb)
	_ = trb.Connect(tra)
	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	_ = tra.Connect(trb)
	_ = trb.Connect(tra)

	assert.Nil(t, tra.SendMessage(trb.addr, []byte("hello")))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, []byte("hello"))
	assert.Equal(t, rpc.From, tra.addr)
}
