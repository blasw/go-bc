package core

import (
	"go-blockchain/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	data := []byte("foo")
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	data := []byte("foo")
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.PublicKey = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}
