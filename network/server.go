package network

import (
	"bytes"
	"fmt"
	"go-blockchain/core"
	"go-blockchain/crypto"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultBlockTime = 5 * time.Second
)

type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProccesor  RPCProccesor
	Transports    []Transport
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
}

type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == 0 {
		opts.BlockTime = DefaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	// if no custom proccessor provided then use server as a default proccessor
	if s.RPCProccesor == nil {
		s.RPCProccesor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}

			if err := s.RPCProccesor.ProccessMessage(msg); err != nil {
				logrus.Error(err)
			}

		case <-s.quitCh:
			break free

		case <-ticker.C:
			if s.isValidator {
				_ = s.createNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) ProccessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.proccessTransaction(t)
	}

	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) proccessTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash":   hash,
			"length": len(s.memPool.Transactions()),
		}).Info("provided tx is already in mempool, skipping")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Info("adding new tx to mempool")

	// TODO: broadcast new tx to peers
	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block...")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
