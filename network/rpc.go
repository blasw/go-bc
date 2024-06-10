package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-blockchain/core"
	"io"

	"github.com/sirupsen/logrus"
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

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, fmt.Errorf("failed to decode message from %s: %w", rpc.From, err)
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message header %v", msg.Header)
	}
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCProccesor interface {
	ProccessMessage(*DecodedMessage) error
}
