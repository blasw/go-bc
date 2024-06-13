package network

import (
	"bytes"
	"go-blockchain/core"
	"go-blockchain/crypto"
	"go-blockchain/types"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/sirupsen/logrus"
)

const (
	DefaultBlockTime = 5 * time.Second
)

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProccesor  RPCProccesor
	Transports    []Transport
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
}

type Server struct {
	ServerOpts
	chain       *core.Blockchain
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == 0 {
		opts.BlockTime = DefaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockChain(genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
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

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				_ = s.Logger.Log("err", err)
			}

			if err := s.RPCProccesor.ProccessMessage(msg); err != nil {
				logrus.Error(err)
			}

		case <-s.quitCh:
			break free

		}
	}

	_ = s.Logger.Log("msg", "Server shutdown")
}

func (s *Server) validatorLoop() {
	_ = s.Logger.Log("msg", "Starting validator loop")

	ticker := time.NewTicker(s.blockTime)

	for {
		<-ticker.C
		_ = s.createNewBlock()
	}
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
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	_ = s.Logger.Log(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempool_length", s.memPool.Len(),
	)

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

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

func (s *Server) createNewBlock() error {
	_ = s.Logger.Log("msg", "creating new block", "cur_height", s.chain.Height())
	curHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	//TODO: delete after testing
	txx := []core.Transaction{}
	for _, tx := range s.memPool.Transactions() {
		txx = append(txx, *tx)
	}

	block, err := core.NewBlockFromPrevHeader(curHeader, txx)
	if err != nil {
		return err
	}

	if block.Sign(*s.PrivateKey) != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.Flush()

	return nil
}

func genesisBlock() *core.Block {
	header := core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}

	return core.NewBlock(&header, nil)
}
