package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeypairSignVerify(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	msg := []byte("hello world!")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pubKey, msg))

	privKey2 := GeneratePrivateKey()
	pubKey2 := privKey2.PublicKey()

	assert.False(t, sig.Verify(pubKey2, msg))

	sig2, err := privKey2.Sign(msg)
	assert.Nil(t, err)

	assert.False(t, sig2.Verify(pubKey, msg))

	assert.False(t, sig.Verify(pubKey, []byte("wrong message")))
}
