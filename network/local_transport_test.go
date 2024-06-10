package network

import (
	"io"
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

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)
	trc := NewLocalTransport("C").(*LocalTransport)

	assert.Nil(t, tra.Connect(trb))
	assert.Nil(t, tra.Connect(trc))

	msg := []byte("hello world")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcc := <-trc.Consume()
	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)
}
