package network

import (
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
	RPCHandler RPCHandler
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime  time.Duration
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

	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if opts.RPCHandler == nil {
		s.ServerOpts.RPCHandler = NewDefaultRPCHandler(s)
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
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
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

func (s *Server) ProccessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash":   hash,
			"from":   from,
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

	return s.memPool.Add(tx)
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
