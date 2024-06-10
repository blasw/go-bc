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

	msg := []byte("hello")

	assert.Nil(t, tra.SendMessage(trb.addr, msg))

	rpc := <-trb.Consume()

	buf := make([]byte, len(msg))
	_, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)

	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tra.addr)
}
