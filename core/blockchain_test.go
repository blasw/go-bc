package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockChain(randomBlock(0))
	assert.NotNil(t, bc)
	assert.Nil(t, err)
	return bc
}

func testAddBlocks(t *testing.T) *Blockchain {
	bc := newBlockchainWithGenesis(t)

	for i := 0; i < 1000; i++ {
		err := bc.AddBlock(randomBlock(uint32(i + 1)))
		assert.Nil(t, err)
	}

	assert.NotNil(t, bc.AddBlock(randomBlock(69)))

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
