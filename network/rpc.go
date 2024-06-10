package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-blockchain/core"
	"io"
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock
)

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	_ = gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type RPCHandler interface {
	HandleRPC(RPC) error //decoder for payload
}

type DefaultRPCHandler struct {
	p RPCProccesor
}

func NewDefaultRPCHandler(p RPCProccesor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		p: p,
	}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return err
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return fmt.Errorf("failed to decode message from %s: %w", rpc.From, err)
		}
		return h.p.ProccessTransaction(rpc.From, tx)
	default:
		return fmt.Errorf("invalid message header %v", msg.Header)
	}
}

type RPCProccesor interface {
	ProccessTransaction(NetAddr, *core.Transaction) error
}
