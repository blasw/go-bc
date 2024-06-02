package core

import (
	"go-blockchain/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(randomBlock(t, 0, types.Hash{}))
	assert.NotNil(t, bc)
	assert.Nil(t, err)
	return bc
}

func testAddBlocks(t *testing.T) *Blockchain {
	bc := newBlockchainWithGenesis(t)

	for i := 0; i < 1000; i++ {
		err := bc.AddBlock(randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1))))
		assert.Nil(t, err)
	}

	assert.NotNil(t, bc.AddBlock(randomBlock(t, 69, types.Hash{})))

	return bc
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Equal(t, len(bc.headers), 1)

	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))

	bc = testAddBlocks(t)

	assert.True(t, bc.HasBlock(1))
	assert.True(t, bc.HasBlock(1000))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000

	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockTooHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, 1, getPrevBlockHash(t, bc, 1))))
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 69, types.Hash{})))
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}
