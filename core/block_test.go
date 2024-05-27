package core

import (
	"bytes"
	"fmt"
	"go-blockchain/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    10,
		Nonce:     912949,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h.Version, hDecode.Version)
	assert.Equal(t, h.PrevBlock, hDecode.PrevBlock)
	assert.Equal(t, h.Timestamp, hDecode.Timestamp)
	assert.Equal(t, h.Height, hDecode.Height)
	assert.Equal(t, h.Nonce, hDecode.Nonce)
}

func TestBlock_Encode_Decode(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     912949,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))

	assert.Equal(t, b, bDecode)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Height:    10,
			Nonce:     912949,
		},
		Transactions: nil,
	}

	h := b.Hash()
	fmt.Println(h)
	assert.False(t, h.IsZero())
}
